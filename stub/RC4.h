//
// KSA (Key Scheduling Algorithm) 
// PRGA (Pseudo-Random Generation Algorithm)
// XOR 
//


#define N 256   // 2^8

void swap(unsigned char *a, unsigned char *b) {
    int tmp = *a;
    *a = *b;
    *b = tmp;
}

void KSA(unsigned char *key,int len,unsigned char *S) {
    int j = 0;
    for(int i = 0; i < N; i++)
        S[i] = i;
    for(int i = 0; i < N; i++) {
        j = (j + S[i] + key[i % len]) % N;
        swap(&S[i], &S[j]);
    }
}

void PRGA(unsigned char *S,unsigned char *plain,int len,unsigned char *cipher) {
    int i = 0;
    int j = 0;

    for(size_t n = 0; n < len; n++) {
        i = (i + 1) % N;
        j = (j + S[i]) % N;
        swap(&S[i], &S[j]);
        int rnd = S[(S[i] + S[j]) % N];
        cipher[n] = rnd ^ plain[n];
    }
}
