// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/wfs

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>
#include <dlfcn.h>

typedef char* (*InitFunc)();
typedef int (*IsInitFunc)();
typedef char* (*AppendFunc)(char* name, unsigned char* data, int dataLen, int compress);
typedef char* (*DeleteFunc)(char* path);
typedef unsigned char* (*GetFunc)(char* path, int* resultLen);
typedef char* (*RenameFunc)(char* path, char* newpath);
typedef void (*FreeMemoryFunc)(void* ptr);
typedef void (*FreeStringFunc)(void* ptr);
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

void* hDll = NULL;

int LoadSharedLib() {
    hDll = dlopen("./wfs.so", RTLD_LAZY);
    if (hDll == NULL) {
        printf("Failed to load wfs.so: %s\n", dlerror());
        return 0;
    }

    Init = (InitFunc)dlsym(hDll, "Init");
    IsInit = (IsInitFunc)dlsym(hDll, "IsInit");
    Append = (AppendFunc)dlsym(hDll, "Append");
    Delete = (DeleteFunc)dlsym(hDll, "Delete");
    Get = (GetFunc)dlsym(hDll, "Get");
    Rename = (RenameFunc)dlsym(hDll, "Rename");
    FreeMemory = (FreeMemoryFunc)dlsym(hDll, "FreeMemory");
    FreeString = (FreeStringFunc)dlsym(hDll, "FreeString");
    Close = (CloseFunc)dlsym(hDll, "Close");
    Has = (HasFunc)dlsym(hDll, "Has");
    GetKeys = (GetKeysFunc)dlsym(hDll, "GetKeys");

    char* error = dlerror();
    if (error) {
        printf("Failed to get function addresses: %s\n", error);
        dlclose(hDll);
        return 0;
    }

    printf("Shared library loaded successfully\n");
    return 1;
}

// Unload shared library
void UnloadSharedLib() {
    if (hDll) {
        dlclose(hDll);
        hDll = NULL;
        printf("Shared library unloaded\n");
    }
}

// Wait for initialization to complete
int WaitForInitialization(int maxWaitSeconds) {
    printf("Waiting for initialization...\n");

    for (int i = 0; i < maxWaitSeconds; i++) {
        sleep(1); // Wait 1 second

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
        printf("Library not initialized\n");
        return;
    }

    if (dataLen == 0 || data == NULL) {
        printf("File does not exist or is empty\n");
        return;
    }

    printf("Data retrieved, length: %d bytes\n", dataLen);
    printf("Hex data: ");
    for (int i = 0; i < dataLen && i < 50; i++) {
        printf("%02X ", data[i]);
    }
    printf("\n");

    printf("Text content: ");
    for (int i = 0; i < dataLen && i < 50; i++) {
        if (data[i] >= 32 && data[i] <= 126) {
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

    printf("\n--- Getting file after delete ---\n");
    dataLen = 0;
    data = Get("renamed_file.txt", &dataLen);
    if (dataLen == 0) {
        printf("File successfully deleted\n");
    }

    printf("\nAll tests completed\n");
}

int main() {
    printf("WFS Shared Library Test Program\n");
    printf("================================\n");

    // 1. Load shared library
    if (!LoadSharedLib()) {
        return 1;
    }

    // 2. Initialize
    printf("Initializing WFS...\n");
    char* initResult = Init("{\"http\": false, \"thrift\": false, \"admin\": false}"); //close all service
    if (initResult) {
        printf("Initialization error: %s\n", initResult);
        FreeMemory(initResult);
        UnloadSharedLib();
        return 1;
    }

    // 3. Wait for initialization
    if (!WaitForInitialization(30)) {
        printf("Initialization failed\n");
        UnloadSharedLib();
        return 1;
    }

    // 4. Run tests
    RunAllTests();

    // 5. Cleanup
    printf("\nCleaning up...\n");
    Close();
    UnloadSharedLib();

    printf("Test program finished\n");
    return 0;
}