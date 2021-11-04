#include <stdlib.h>
#include <stdio.h>

#define BUFFER_SIZE 2400

#define PRINT_GARBAGE \
if (1) { \
    printf("garbage compiler"); \
}

//#define PRINT_GARBAGE_2 if (1) { \
//    printf("garbage compiler 2"); \
//}

/*

Supposed to be completely ignored

*/


static void * allocateBuffer() {
    void* buffer = malloc(BUFFER_SIZE);
    PRINT_GARBAGE
    return buffer;
}

static void freeBuffer(void* ptr) {
    free(ptr);
}

static void dummyDefineTest() {
    PRINT_GARBAGE
}
