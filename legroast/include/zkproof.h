#ifndef ZKPROOF_H
#define ZKPROOF_H

#include <openssl/rand.h>
#include <string.h>
#include <stdint.h>
#include <stdio.h>
#include "parameters.h"
#include "merkletree.h"

#define SK_BYTES SEED_BYTES

//the order of the shares in memory
#define SHARE_K 0
#define SHARES_TRIPLE (SHARE_K + 1)
#define SHARES_R (SHARES_TRIPLE + 3)
#define SHARES_PER_PARTY (SHARES_R + RESIDUOSITY_SYMBOLS_PER_ROUND)

typedef __int128 int128_t;
typedef unsigned __int128 uint128_t;

typedef struct 
{
    unsigned char seed_trees[ROUNDS][SEED_BYTES*(2*PARTIES-1)];
    uint128_t shares[ROUNDS][PARTIES][SHARES_PER_PARTY];
    uint128_t sums[ROUNDS][SHARES_PER_PARTY];
    uint128_t indices[ROUNDS*RESIDUOSITY_SYMBOLS_PER_ROUND];
} prover_state;

#define MESSAGE1_COMMITMENT(mes) (mes)
#define MESSAGE1_DELTA_K(mes) (MESSAGE1_COMMITMENT(mes) + HASH_BYTES)
#define MESSAGE1_DELTA_TRIPLE(mes) (MESSAGE1_DELTA_K(mes) + ROUNDS*sizeof(uint128_t))
#define MESSAGE1_BYTES (MESSAGE1_DELTA_TRIPLE(0) + ROUNDS*PRIME_BYTES)

#define CHALLENGE1_BYTES (ROUNDS*RESIDUOSITY_SYMBOLS_PER_ROUND*sizeof(uint32_t))

#define MESSAGE2_OUTPUT(mes) (mes)
#define MESSAGE2_BYTES (MESSAGE2_OUTPUT(0) + ROUNDS*RESIDUOSITY_SYMBOLS_PER_ROUND*PRIME_BYTES)

#define CHALLENGE2_EPSILON(ch) (ch)
#define CHALLENGE2_LAMBDA(ch) (CHALLENGE2_EPSILON(ch) + ROUNDS*PRIME_BYTES)
#define CHALLENGE2_BYTES (CHALLENGE2_LAMBDA(0) + ROUNDS*RESIDUOSITY_SYMBOLS_PER_ROUND*PRIME_BYTES)

#define MESSAGE3_HASH(mes) (mes)
#define MESSAGE3_ALPHA(mes) (MESSAGE3_HASH(mes) + HASH_BYTES)
#define MESSAGE3_BETA(mes) (MESSAGE3_ALPHA(mes) + ROUNDS*PRIME_BYTES)
#define MESSAGE3_BYTES (MESSAGE3_BETA(0) + ROUNDS*PRIME_BYTES)

#define CHALLENGE3_BYTES (sizeof(uint32_t[ROUNDS]))

#define MESSAGE4_SEEDS(mes) (mes)
#define MESSAGE4_COMMITMENT(mes) (MESSAGE4_SEEDS(mes) + ROUNDS*PARTY_DEPTH*SEED_BYTES)
#define MESSAGE4_BYTES (MESSAGE4_COMMITMENT(0) + ROUNDS*HASH_BYTES)

#define RESPONSE1_BYTES 1
#define RESPONSE2_BYTES 1

void keygen(unsigned char *pk, unsigned char *sk);
void commit(const unsigned char *sk, const unsigned char *pk, unsigned char *message1, prover_state *state);
void generate_challenge1(const unsigned char *hash, unsigned char *challenge1 );
void respond1(const unsigned char *sk, const unsigned char *pk, const unsigned char *challenge1, unsigned char *message2, prover_state *state);
void generate_challenge2(const unsigned char *hash, unsigned char *challenge2 );
void respond2(const unsigned char *sk, const unsigned char *pk, const unsigned char *challenge1, const unsigned char *challenge2, unsigned char *message2, unsigned char *message3, prover_state *state);
void generate_challenge3(const unsigned char *hash, unsigned char *challenge3 );
void respond3(const unsigned char *sk, const unsigned char *pk, const unsigned char *challenge1, const unsigned char *challenge2, const unsigned char *challenge3, unsigned char *message4, prover_state *state);
int check(const unsigned char *pk, const unsigned char *message1, const unsigned char *challenge1, const unsigned char *message2, const unsigned char *challenge2, const unsigned char *message3, const unsigned char *challenge3, const unsigned char *message4);

unsigned char legendre_symbol(uint128_t *a);

#endif