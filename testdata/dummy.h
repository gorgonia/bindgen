struct dummy;
typedef struct dummy dummy_t;

struct dummy2;
typedef struct dummy2 *dummy2_t;

typedef enum err_enum {
	SUCCESS = 0,
	FAILURE = 1,
} error;

typedef enum fntype {
	DOES_SOMETHING = 0,
	DOES_NOTHING = 1,
} fntype_t;

typedef int foo; 
typedef int context;

void func1i(int* a);
void func1f(foo a); 
void func1fp(foo* a);

void func2i(int a, int b);
void func2f(foo a, int b);

error funcErr(const int* a, foo* retVal);
error funcCtx(const context* ctx, foo a, foo* retVal);