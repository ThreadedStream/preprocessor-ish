#include <stdlib.h>
#include <stdio.h>

#define BUFFER_SIZE 2048

#define GO_SILLY_V2 if (0) \
printf("you're gonna be out of luck")


#define GO_SILLY \
if (1) \
    printf("went just completely silly");


int main(int argc, const char* argv[]) {

    void* ptr = malloc(BUFFER_SIZE);

    GO_SILLY

    GO_SILLY_V2

    free(ptr);

    return 0;
}

