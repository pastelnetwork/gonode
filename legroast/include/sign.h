#ifndef SIGN_H
#define SIGN_H 

#include <stdlib.h>

void *aligned_alloc( size_t alignment, size_t size );

#include "parameters.h"
#include "zkproof.h"

#define SIG_MESSAGE1(sig) (sig)
#define SIG_MESSAGE2(sig) (SIG_MESSAGE1(sig) + MESSAGE1_BYTES)
#define SIG_MESSAGE3(sig) (SIG_MESSAGE2(sig) + MESSAGE2_BYTES)
#define SIG_MESSAGE4(sig) (SIG_MESSAGE3(sig) + MESSAGE3_BYTES)
#define SIG_BYTES (SIG_MESSAGE4(0) + MESSAGE4_BYTES) 

void sign(const unsigned char *sk,  const unsigned char *pk,const unsigned char *m, uint64_t mlen, unsigned char *sig, uint64_t *sig_len);
int verify(const unsigned char *pk, const unsigned char *m, uint64_t mlen, const unsigned char *sig); 

#endif