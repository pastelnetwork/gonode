#ifndef PARAMETERS_H
#define PARAMETERS_H 

//#define FAST
#define MIDDLE
//#define COMPACT

#define LEGENDRE
//#define POWER


#define PRIME_BYTES 16
#define SEED_BYTES 16
#define HASH_BYTES 32
#define PK_DEPTH 15

#ifdef LEGENDRE

	#ifdef FAST
		#define ROUNDS 54 
		#define RESIDUOSITY_SYMBOLS_PER_ROUND 9 
		#define PARTY_DEPTH 4 
	#endif

	#ifdef MIDDLE
		#define ROUNDS 37
		#define RESIDUOSITY_SYMBOLS_PER_ROUND 12 
		#define PARTY_DEPTH 6 
	#endif

	#ifdef COMPACT
		#define ROUNDS 26
		#define RESIDUOSITY_SYMBOLS_PER_ROUND 16 
		#define PARTY_DEPTH 8 
	#endif

#endif

#ifdef POWER

	#ifdef FAST
		#define ROUNDS 39 
		#define RESIDUOSITY_SYMBOLS_PER_ROUND 4 
		#define PARTY_DEPTH 4 
	#endif

	#ifdef MIDDLE
		#define ROUNDS 27
		#define RESIDUOSITY_SYMBOLS_PER_ROUND 5 
		#define PARTY_DEPTH 6 
	#endif

	#ifdef COMPACT
		#define ROUNDS 21
		#define RESIDUOSITY_SYMBOLS_PER_ROUND 5
		#define PARTY_DEPTH 8 
	#endif

#endif

#define PARTIES (1<<PARTY_DEPTH)
#define PK_BYTES (1<<(PK_DEPTH-3))

#define LEAVES (1 << DEPTH)
#define LEAF_BYTES (A_COLS*sizeof(uint16_t))
#define TREE_BYTES ((2*LEAVES-1)*HASH_BYTES)
#define PATH_BYTES (DEPTH*HASH_BYTES)

#include "libkeccak.a.headers/SimpleFIPS202.h"
#define HASH(data,len,out) SHAKE128(out, HASH_BYTES, data, len);
#define EXPAND(data,len,out,outlen) SHAKE128(out, outlen, data, len);

#endif
