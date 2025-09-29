// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/wfs

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <windows.h>

typedef char* (*InitFunc)();
typedef int (*IsInitFunc)();
typedef char* (*AppendFunc)(char* name, unsigned char* data, int dataLen, int compress);
typedef char* (*DeleteFunc)(char* path);
typedef unsigned char* (*GetFunc)(char* path, int* resultLen);
typedef char* (*RenameFunc)(char* path, char* newpath);
typedef void (*FreeMemoryFunc)(void* ptr);
typedef void (*FreeStringFunc)(char* s);
typedef void (*CloseFunc)();
typedef int (*HasFunc)(char* path);
typedef char* (*GetKeysFunc)(long long fromId, int limit);

InitFunc Init = NULL;
IsInitFunc IsInit = NULL;
AppendFunc Append = NULL;
DeleteFunc Delete = NULL;
GetFunc Get = NULL;
RenameFunc Rename = NULL;
FreeMemoryFunc FreeMemory = NULL;
FreeStringFunc FreeString = NULL;
CloseFunc Close = NULL;
HasFunc Has = NULL;
GetKeysFunc GetKeys = NULL;

HMODULE hDll = NULL;

int LoadDLL() {
    hDll = LoadLibraryA("wfs.dll");
    if (hDll == NULL) {
        printf("Failed to load wfs.dll\n");
        return 0;
    }

    Init = (InitFunc)GetProcAddress(hDll, "Init");
    IsInit = (IsInitFunc)GetProcAddress(hDll, "IsInit");
    Append = (AppendFunc)GetProcAddress(hDll, "Append");
    Delete = (DeleteFunc)GetProcAddress(hDll, "Delete");
    Get = (GetFunc)GetProcAddress(hDll, "Get");
    Rename = (RenameFunc)GetProcAddress(hDll, "Rename");
    FreeMemory = (FreeMemoryFunc)GetProcAddress(hDll, "FreeMemory");
    FreeString = (FreeStringFunc)GetProcAddress(hDll, "FreeString");
    Close = (CloseFunc)GetProcAddress(hDll, "Close");
    Has = (HasFunc)GetProcAddress(hDll, "Has");
    GetKeys = (GetKeysFunc)GetProcAddress(hDll, "GetKeys");

    if (!Init || !IsInit || !Append || !Delete || !Get || !Rename ||
        !FreeMemory ||!FreeString || !Close || !Has || !GetKeys) {
        printf("Failed to get all function addresses\n");
        FreeLibrary(hDll);
        return 0;
    }

    printf("DLL loaded successfully\n");
    return 1;
}

void UnloadDLL() {
    if (hDll) {
        FreeLibrary(hDll);
        hDll = NULL;
        printf("DLL unloaded\n");
    }
}

// Wait for initialization to complete
int WaitForInitialization(int maxWaitSeconds) {
    printf("Waiting for initialization...\n");

    for (int i = 0; i < maxWaitSeconds; i++) {
        Sleep(1000); // Wait 1 second

        int status = IsInit();
        if (status == 1) {
            printf("Initialization completed\n");
            return 1;
        }

        printf("Waiting... (%d/%d seconds)\n", i + 1, maxWaitSeconds);
    }

    printf("Initialization timeout\n");
    return 0;
}

// Test Append function
void TestAppend() {
    printf("\n=== Testing Append Function ===\n");

    // Prepare test data
    unsigned char testData[] = {0x48, 0x65, 0x6C, 0x6C, 0x6F, 0x20, 0x57, 0x46, 0x53}; // "Hello WFS"
    char* result = Append("test_file.txt", testData, sizeof(testData), 0);

    if (result) {
        printf("Append failed: %s\n", result);
        FreeMemory(result);
    } else {
        printf("Append successful\n");
    }
}

// Test Get function
void TestGet() {
    printf("\n=== Testing Get Function ===\n");

    int dataLen = 0;
    unsigned char* data = Get("test_file.txt", &dataLen);

    if (dataLen == -1) {
        printf("DLL not initialized\n");
        return;
    }

    if (dataLen == 0 || data == NULL) {
        printf("File does not exist or is empty\n");
        return;
    }

    printf("Data retrieved, length: %d bytes\n", dataLen);
    printf("Hex data: ");
    for (int i = 0; i < dataLen && i < 50; i++) { // Show max 50 bytes
        printf("%02X ", data[i]);
    }
    printf("\n");

    // Try to display as text (if it's text data)
    printf("Text content: ");
    for (int i = 0; i < dataLen && i < 50; i++) {
        if (data[i] >= 32 && data[i] <= 126) { // Printable characters
            printf("%c", data[i]);
        } else {
            printf(".");
        }
    }
    printf("\n");

    FreeMemory(data);
}

// Test Rename function
void TestRename() {
    printf("\n=== Testing Rename Function ===\n");

    char* result = Rename("test_file.txt", "renamed_file.txt");

    if (result) {
        printf("Rename failed: %s\n", result);
        FreeMemory(result);
    } else {
        printf("Rename successful\n");
    }
}

// Test Delete function
void TestDelete() {
    printf("\n=== Testing Delete Function ===\n");

    char* result = Delete("renamed_file.txt");

    if (result) {
        printf("Delete failed: %s\n", result);
        FreeMemory(result);
    } else {
        printf("Delete successful\n");
    }
}

// Test Has function
void TestHas() {
    printf("\n=== Testing Has Function ===\n");

    int exists = Has("renamed_file.txt");

    if (exists == -1) {
        printf("DLL not initialized\n");
    } else if (exists == 1) {
        printf("File exists\n");
    } else {
        printf("File does not exist\n");
    }
}

// Test GetKeys function
void TestGetKeys() {
    printf("\n=== Testing GetKeys Function ===\n");

    char* keysJSON = GetKeys(1, 100); // fromId=1, limit=10
    printf("Keys: %s\n", keysJSON);
    FreeMemory(keysJSON);
}

// Run all tests
void RunAllTests() {
    printf("Starting all tests...\n");


    TestAppend();
    TestGet();
    TestRename();
    TestHas();
    TestGetKeys();

    // Try to get after rename
    printf("\n--- Getting file after rename ---\n");
    int dataLen = 0;
    unsigned char* data = Get("renamed_file.txt", &dataLen);
    if (dataLen > 0) {
        printf("Renamed file exists, length: %d bytes\n", dataLen);
        FreeMemory(data);
    } else {
        printf("Renamed file does not exist\n");
    }



    // TestDelete();

    // Try to get after delete
    // printf("\n--- Getting file after delete ---\n");
    //    dataLen = 0;
    //    data = Get("renamed_file.txt", &dataLen);
    //    if (dataLen == 0) {
    //        printf("File successfully deleted\n");
    //    }

    printf("\nAll tests completed\n");
}

int main() {
    printf("WFS DLL Test Program\n");
    printf("====================\n");

    // 1. Load DLL
    if (!LoadDLL()) {
        return 1;
    }

    // 2. Initialize
    printf("Initializing WFS...\n");
    char* initResult = Init("{\"http\": false, \"thrift\": false, \"admin\": false}");
    if (initResult) {
        printf("Initialization error: %s\n", initResult);
        FreeMemory(initResult);
        UnloadDLL();
        return 1;
    }

    // 3. Wait for initialization
    if (!WaitForInitialization(30)) {
        printf("Initialization failed\n");
        UnloadDLL();
        return 1;
    }

    // 4. Run tests
    RunAllTests();

    // 5. Cleanup
    printf("\nCleaning up...\n");
    Close();
    UnloadDLL();

    printf("Test program finished\n");
    return 0;
}