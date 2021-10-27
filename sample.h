#include <stdlib.h>
#include <stdio.h>

#define BUFFER_SIZE 12400
#define PRINT_GARBAGE \
if (1) { \
    printf("garbage compiler"); \
}

static void * allocateBuffer() {
    void* buffer = malloc(BUFFER_SIZE);
    PRINT_GARBAGE
    return buffer;
}

