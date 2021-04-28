#ifndef MERKLETREE_H
#define MERKLETREE_H

#include <stdio.h>
#include "stdint.h"
#include <string.h>

#include "parameters.h"

void generate_seed_tree(unsigned char *seed_tree);

void release_seeds(unsigned char *tree, uint32_t unopened_index, unsigned char *out);
void fill_down(unsigned char *tree, uint32_t unopened_index, const unsigned char *in);
 
/*

void build_tree(const unsigned char *data,int leaf_size, int depth, unsigned char *tree);

void get_path(const unsigned char *tree, int depth, int leaf_index, unsigned char *path);
void follow_path(const unsigned char *leaf, int leaf_size, int depth, const unsigned char *path, int index, unsigned char *root);

void print_seed(const unsigned char *seed);
void print_hash(const unsigned char *hash);
void print_tree(const unsigned char *tree, int depth);

void hash_up(unsigned char *data, unsigned char *indices, const unsigned char *in, int in_len, unsigned char *root); */

#endif