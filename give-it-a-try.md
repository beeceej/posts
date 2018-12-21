<!--id: 3-->
<!--title: Give it a Try! -->
<!--author: Brian Jones-->
<!--postedAt: July 23rd, 2018-->
<!--updatedAt: August 7th, 2018-->
<!--visible: true-->

5 years ago if you told me i'd be hacking in LISP, I would have told you, \"You're crazy\". I've heard too many people talk about how bad it is. Most arguments revolving around functional programming not being practical. Or about all of the parentheses (_**How do you even match those up?!?!**_). I agree functional programming can seem relatively opaque, especially if that person is schooled in OOP languages, or the only exposure to functional programming emphasizes on theory, as opposed to practical application. But, it is extremely short sighted to disregard a language on syntax alone. We should spend half the time on learning the theory, and half the time on practical applications. If we didn't hold misconceptions about languages and features maybe we'd all be further along? Maybe we'd stop arguing about inconsequential matters.

Hindsight is 20/20, so now I'm taking time to look into the once alien lands of programming. There is no silver bullet, and alien concepts definitely translate to being more productive (This is probably why people fall so quickly into their comfort zones). It's very easy to delve into Java or Python, or whatever and only gloss over some of the more exotic CS topics; as a result, many only skim the surface. If you only know Java, then you are unable to realize how verbose the syntax actually is. similarly, it's very difficult to appreciate the areas that it excels in, IDE's, and _mostly_ safe types. Learning new languages and concepts, even if you don't use it to be more productive, by default enables you to be more effective in your language of choice. For someone going from `LANG_X -> LANG_Y` it could be a scary jump. And rightfully so, new paradigms mean new problems which is scary. But remember, new problems usually means new successes, it's all about tradeoffs.

I'd argue after a couple weeks, to months, if one drops preconceptions they can start to see the language objectively for its pros and cons. Context is key to learning for oneself.

- The Java developer starts to see how not declaring types on _**every**_ variable definition saves time.

- The Python developer starts to see how types aren't just a nuisance, but they provide structure in your code that _**eliminates a multitude of errors which would be caught at run-time in python**_.

- The OOP programmer begins to see the benefits of immutable data and referential transparency

- The (Pure) Functional Programmer begins to see that interacting with the world is much easier when you can make side-effects, and one can be just as productive when the side affects aren't enforced by the type system

What arises from all of these ideas is what becomes a set of best-practices (either enforced by the language, or not). It's easy to blindly follow best practices. Though, i'd argue in the long run for a developer, it is more beneficial to question best practices, and use learnings from experience (where necessary). Standards are built over time, standards also change over time. With time standards age. Depending on the project at hand you'll have to make a decision to be productive, or to expand your horizons... Anyway, playing with Common LISP has been fun... And here's an observation I've come across...

---

## Common LISP with parens. (they're really not bad)

```lisp
(defun high-or-low (guess target)
 (when (> guess target) (format t \"**A bit High**~%\"))
 (when (< guess target) (format t \"**A bit Low**~%\")))

(defun is-winner? (b answer)
  (if b (print \"**You Win!**\") (format t \"~%**You Lose!**~%~%**The answer was ~a**\" answer)))

(defun capture-user-guess ()
  (handler-case
    (parse-integer (read-line))
    (t (c)
      (declare (ignore c))
      (format t \"**Invalid input, expecting a number**~%\")
      (capture-user-guess ))))

(defun .guessing-game (curr-try max-tries target)
  (if (>= curr-try max-tries)
    nil
    (progn
      (format t \"**Guess ~a**~%\" curr-try)
      (let ((guess (capture-user-guess )))
        (if (eq guess target)
          t
          (progn
            (high-or-low guess target)
            (.guessing-game (+ 1 curr-try) max-tries target)))))))

(defun guessing-game (num-range tries target)
  (format t \"**Welcome to the guessing game, Would you like to play? (Y/N)**~%\")
  (let ((answer (read-line)))
    (if (string-equal answer \"Y\")
      (progn
        (format t \"**Guess a number 0 - ~a** ~%~%\" num-range)
        (is-winner? (.guessing-game 0 tries target) target))
      (format t \"**Cya**~%\"))))
```

## Common Lisp without Parens, looks a lot like python, right?!

```python
defun high-or-low guess target
 when > guess target format t \"**A bit High**~%\"
 when < guess target format t \"**A bit Low**~%\"

defun is-winner? b answer
  if b print \"**You Win!**\" format t \"~%**You Lose!**~%~%**The answer was ~a**\" answer

defun capture-user-guess
  handler-case
    parse-integer read-line
    t c
      declare ignore c
      format t \"**Invalid input, expecting a number**~%\"
      capture-user-guess

defun .guessing-game curr-try max-tries target
  if >= curr-try max-tries
    nil
    progn
      format t \"**Guess ~a**~%\" curr-try
      let guess capture-user-guess
        if eq guess target
          t
          progn
            high-or-low guess target
            .guessing-game + 1 curr-try max-tries target

defun guessing-game num-range tries target
  format t \"**Welcome to the guessing game, Would you like to play? Y/N**~%\"
  let answer read-line
    if string-equal answer \"Y\"
      progn
        format t \"**Guess a number 0 - ~a** ~%~%\" num-range
        is-winner? .guessing-game 0 tries target target
      format t \"**Cya**~%\"
```

- Doesn't this lisp code look like super functional python?

[Source for above code](https://github.com/beeceej/guessing-game)
