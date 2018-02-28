# bindgen
package bindgen is a package that aids in the generation of bindings and idiomatic Go interfaces to C-libraries.

As it exists it's a collection of utility data structures and functions built on top of [cznic/cc](https://github.com/cznic/cc) to simplify usage of the C parser. The translation has to be written manually still. Each package would have different translation needs, hence the minimal package.

Generation of bindings and interfaces is primarily an exercise syntactic parsing and translation. The semantics of the C language isn't as important, given the Go `C` pseudopackage obeys all the calling semantics of C anyway. 

Much of the architecure was drawn directly from the gonum BLAS generation package. 

# Examples

The running example in this package will involve parsing a dummy C header that is defined like this:

```
typedef enum err_enum {
	SUCCESS = 0,
	FAILURE = 1,
} error;

typedef int foo; 
typedef int context;

void func1i(int* a);
void func1f(foo a); 
void func1fp(foo* a);

void func2i(int a, int b);
void func2f(foo a, int b);

error funcErr(const int* a, foo* retVal);
error funcCtx(const context* ctx, foo a, foo* retVal);
```

# Pseudo-Interactive Conversion

This package provides pseudo-interactive conversions. It's NOT a [wizard](https://en.wikipedia.org/wiki/Wizard_(software)) (remember those?), but provides wizard-like functionalities. The key idea is to have a main package that keeps changing, essentially treating this like a REPL.

Here's a sample process used. We'll start with a main package. Our directory looks like this:

```
.
└── main.go
```

First, use `Explore()` to find out if the header can be easily parsed. 

After you're satisfied, use the `GenIgnored` and `GenNameMap` functions. These functions generate strings that are valid Go code. Write them to a file, which we will call `mappings.go` like so:

```
.
├── main.go
└── mappings.go
```

Now that the mappings are generated, you can now safely modify your  main package to no longer generate the mappings. Instead, you can now work on the translation parts.


# vs c-for-go

Should you use this or [c-for-go](https://github.com/xlab/c-for-go)? This package is not as full-featured. If you want an automatic solution, use c-for-go. Unfortunately, this author has found the results are not as idiomatic.

Consider the  header in the example above again.

These are the idiomatic equivalents in Go of the last two:

```
// error funcErr(const int* a, foo* retVal);
func funcErr(a int) (retVal *foo, err error) {} 

// error funcCtx(const context* ctx, foo a, foo* retVal);
func (ctx *context) funcCtx(a foo) (retVal *foo, err error){}
```

With c-for-go, translating the above to those would entail a very complicated translation portion YAML file. A few attempts were made to use c-for-go to translate CUDA bindings to Go, but none were satisfactorily idiomatic. 

But we're programmers, so why not program? This is the basis for this package - to provide easy-to-access features necessary to do the translation. 

Granted, how you'd want to structure your code as "idiomatic" is up to you. Since this package is used quite a bit in generating the [Go bindings for CUDA](https://github.com/gorgonia/cu), there are a lot of repeated patterns like those two examples above.

The CUDA libraries also come with certain patterns that make translation easier - things like `const` parameters makes things clearer, and simple rules can be written in code to translate them. Sure, you end up not having something easily defined as YAML files, but you also end up with nicer (imo) code.

**Conclusion**: If you want something automatic, use the excellent (and well documented) [c-for-go](https://github.com/xlab/c-for-go). If you want more control over the translation process, use this package.


# Tips For Writing Bindings

Sometimes your includes directories can be a little wonky. Since the task at hand is writing/generating the bindings and Go API for said bindings, these are a few things I've found helpful:

1. `gcc -E -P FILE.h > CLEAN.h` This expands all the macros that some files have. 
2. Mock data structures: some types that come from `#include` are complicated chains of even more `#include`s. Just replace them with dummy data structures.
3. The same for `typedef` that are too complicated and involve too many other files. 




# Licence 

This package is MIT licenced. It is derived from the Gonum generator package which is licenced under the _____(TODO). 
