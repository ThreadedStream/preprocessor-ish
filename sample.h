#include <stdlib.h>
#include <stdio.h>

#define BUFFER_SIZE 12400
#define PRINT_GARBAGE \
if (1) { \
    printf("garbage compiler"); \
}

#define PRINT_GARBAGE_2 if (1) { \
    printf("garbage compiler 2"); \
}


static void * allocateBuffer() {
    void* buffer = malloc(BUFFER_SIZE);
    PRINT_GARBAGE
    return buffer;
}

static void freeBuffer(void* ptr) {
    free(ptr);
}