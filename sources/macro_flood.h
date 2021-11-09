#include <stdlib.h>
#include <stdio.h>

#define BUFFER_SIZE 2048

#define GO_SILLY \
if (1) \
    printf("went just completely silly");


int main(int argc, const char* argv[]) {

    void* ptr = malloc(BUFFER_SIZE);

    GO_SILLY

    free(ptr);

    return 0;
}

