<!--id: 9-->
<!--title: Two Weeks With Emacs-->
<!--author: Brian Jones-->
<!--visible: true-->

The tools a developer use on a daily basis can be near and dear to his or her heart. Much like a warrior with his sword and shield, programmers brandish a language of choice with skill. If programming languages, or paradiagms are the sword, then a programmers chosen editor is like the trusty shield (I don't know... the analogy is kind of falling apart here, but we must push on). The shield, err editor helps the programmer every day through things like **auto completion**, **syntax-highlighting**, **error checking**, **extensability**, **configurability** **language integration**, **portability**, **beauty** and so many more... Like every warrior, every programmer will have different tastes. Some warriors want to go into battle against enemies (business problems) and get the work done then go home.

With all that in mind, I decided I wanted to try out emacs. I'd heard people shout EMACS, and VIM and EMACS vs VIM. I'd already casually used Vim for simple file editing, but I had't yet seriously used emacs, and I've been getting into lisp lately; so I thought it'd be a good time to sharpen my elisp sword and wield my emacs shield. I don't think I thought I'd make a permanent switch from VS Code to emacs, but after two weeks, I can definitely see that happening. All of this with the caveat that I'm not using pure emacs... I'm using spacemacs with EVIL mode (a VI Layer on top of the emacs editor, did someone say **Extensible?!**). This actually brings me to the first order of business, emacs really is extensible. Every time I've gone searching for a plugin, or plugin equivalent, I haven't been disapointed. The spacemacs layer system is awesome, even if you do have to go to the develop branch to pull in functionality from time to time. You want terraform support? you got it there's [https://github.com/syohex/emacs-terraform-mode](https://github.com/syohex/emacs-terraform-mode) and [https://github.com/rafalcieslak/emacs-company-terraform](https://github.com/rafalcieslak/emacs-company-terraform). If you want Golang support, there's [https://github.com/syl20bnr/spacemacs/tree/master/layers/%2Blang/go](https://github.com/syl20bnr/spacemacs/tree/master/layers/%2Blang/go), there's python, haskell, rust and of course the first class common-lisp support with SLIME. You really can't go wrong, and these major modes aren't difficult to configure either. In most cases they just work. Sometimes they take a little bit of fiddling to get things to work, but it hasn't been too bad. So, from the start, the first thing I noticed was how easily configurable, and customizable emacs is. I had to go down a couple rabbit holes though, for example learning pyenv, but I think that was well worth it, even if it was frustrating because I was also learning VIM, and Emacs at the same time.

I won't lie, everything wasn't roses; I was still getting comfortable with modal editing, and I for sure had a drop in productivity. But, after day 3 or 4 I started gaining some muscle memory and began breaking even as far as text-editing speed. On the emacs/spacemacs side, there was also a learning curve. Mostly around how to manage different projects and buffers. Once I figured out how to use [projectile](https://github.com/bbatsov/projectile) all of that went away. This is I think is one of the largest ways emacs as improved my workflow. To give an example: 

I was working in a new python project and was faced with a problem I knew I had solved before, but couldn't remember exactly how to solve it, but emacs allowed me to switch project layouts by pressing `SPACE l #key` then `SPACE s g p` to search for a file containing the snippet of code. I added all of that to my copy buffer, then switched back to my original project by `SPACE l l #keyoflastproject` and pasted the example right into the file with an example of how to solve the problem. It took longer to write this explanation than to actually do that. If I were using VSCode, or another editor I probably would have went to my terminal, `cd`'d into the other project directory, Opened up the project in vscode by `code .` then executed a global project search for the code snippet, and then went back to the other project with the example in hand.

The key difference in the two different workflows, is that in emacs, I never had to leave emacs (haha I know, running joke, put the kitchen sink in emacs). But that is Awesome for productivity, emacs allows me execute MORE of my development workflow in only one program. This means less context switches, which means increased productivity. The crazy thing is, the way I approached that problem probably isn't even as efficient as it could have been, I'm sure I'll learn a better way of handling such situations in the future.

The largest hurdle early on has been just getting acclimated to approaching learning curves. Developers do generally well with learning curves, but eventually you hit a point where you are satisfied with your level of productivity, and sometimes it doesn't feel like it's worth the time to face another challenge. I've had to force myself to be OK with degraded productivity, and feeling useless at something atleast for a couple days (Sorry employer) for a bit of productivity gain. And now, my reward is living on the fringes of software development, with a tool that has stood the test of time. My sword and shield is now so versatile that I can modify them at will, and that is true productivity.

I Apologize if anyone here was looking for some crazy emacs/spacemacs tips and tricks, but I'm still learning, and honestly don't have anything revolutionary. But, here are a list of concepts which may help you become more productive faster than me.


(These are spacemacs specific, sorry pure emacs'ers)

## Projectile
Project management in emacs made awesome.  
`SPACE p p`  Allows you to search through all of your projects. And opens a project in the current layer.  
`SPACE p l`  Allows you to search through all of your projects. And opens a project in a new layer. Once you do this once, you can use `SPACE l #keyofyourproject` to quickly switch projects.  
`SPACE p b`  View all buffers in the current projectile project.  

## Windows  
How your buffers are visually displayed  
`SPACE w -`  Splits your window in half horizontally  
`SPACE w /`  Splits your window in half vertically  
`SPACE w w`  Shifts focus to the next window  
`SPACE w TAB`  Shifts focus to the previous window  
`:q` removes the current window, if there is only one window up, then quits emacs (whoops)  

## Buffers  
Decoupled from windows, buffers can be placed into any window  
`SPACE b b`  see a list of all buffers  
`SPACE b d`  kill a buffer,  
`SPACE TAB`  Switch the current window pane to the previous buffer  


## Searching  
`SPACE s g p`  search through a projectile project with grep 


