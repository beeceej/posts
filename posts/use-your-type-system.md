<!--id: 9-->
<!--title: Use Your Type System-->
<!--author: Brian Jones-->
<!--visible: true-->

Have you ever seen Java code with weakly typed parameters in its constructor? It looks something like this...


```java
public class Dog

  public Dog(String breed, String coatType) {
    // Initialization Code here with potential validity checks, is the breed a valid breed, is the coat type even a coat type? did I mix the parameters up? etc...
  }
}
```

All is fine until someone tries to instantiate a dog like

```java
new Dog("Long haired, brown", "Australian Shepherd");
```

The above is just plain wrong, and it could have been prevented by making use of the type system however limiting it may be (Haskell people don't yell at me!!). A better piece of code will look like:

```java
public class Dog {
  public Dog(Breed breed, Coat coatType) {
    // Initialization Code here with no initialization checks, to even get a Breed it must be built correctly, same with Coat
  }
}
```

And you can create a dog by:


```java
final Optional<Breed> maybeBreed = Breed.of("Australian Shepherd")
final Optional<Coat> maybeCoat = Coat.of("Thick", "Brown", "Spotted")
final BiFunction<Optional<Breed>, Optional<Coat>, Optional<Dog>> makeDog = (maybeBreed, maybeCoat) -> {
  return maybeBreed.flatMap(breed -> {
    maybeCoat.flatMap(coat -> {
      Optional.of(new Dog(breed, coat));
    })
  })
};
final Optional<Dog> maybeDog = makeDog.apply(maybeBreed, maybeCoat);
// Now go outside and play with the dog you just made!
```


With the above code, you can be sure that your dog was created, and because Breed and Coat are **_using your type system_** you can be sure your dog wasn't created with cat features or something crazy like that, if it were, then you wouldn't have a value in your Optional. Note, if you wanted to know WHY creating your dog failed, you could swap out the `Optional's` with an Either type. Further, you could attach useful predicates to your Breed, and Coat class like this:

```java
public class Breed {

  public final static Predicate<Breed> isLarge= breed -> {
    // However you determine a large breed of dog
  };
}
```
then in code it's easy to ask your object questions like:

```java
final Breed aBreed;

Breed.isLarge.test(aBreed); // returns true or false.
```


I'm writing java examples because everyone hates Java's type system... No one said it had to be pretty or elegant, but it gets the job done.
