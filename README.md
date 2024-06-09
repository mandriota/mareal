# Mareal
The high-level scripting Mareal programming language.

## Installation
```sh
go install github.com/mandriota/mareal
```

## Usage
```sh
mareal file.mr
```

## Example: Fibonacci Numbers
```lisp
_

(def fib-helper (lambda (quote n x y)
   (if n (fib-helper (- n 1) (+ x y) x) x)))

(def fib (lambda (quote n)
   (fib-helper n 0 1)))

(put (_ "fibonacci: "
     (fib (num
       (get "enter number: ")))) nl)
```

You can find more examples in ./examples directory.
