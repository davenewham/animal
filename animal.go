package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

var reader = bufio.NewReader(os.Stdin) //we'll be needing this

func main() {
	playIntroMessage()

	run := true
	for run {
		fmt.Println("Are you thinking of an animal?")

		switch strings.ToLower(getInput())[:1] {
		case "n": //No
			run = false //Quit the game
		case "l": //List
			listKnownAnimals()
		case "y": //Yes
			askQuestions()
		default:
			fmt.Println("I don't recognise that command")
		}
	}
}

func playIntroMessage() {
	fmt.Println("Play 'Guess the Animal'")
	fmt.Println("Think of an animal and the computer will try to guess it...")
}

func listKnownAnimals() {
	fmt.Println("Animals I already know are:")
	for _, v := range animals {
		fmt.Println(v)
	}
}

func askQuestions() {
	//I guess we could do this properly recursively but meh

	qIndex := 0
	t := askQuestion(qIndex) //Ask the first question

	//ask a question or guess an animal
	for t.Type == QUESTION {
		qIndex = t.Index //store this, in case t comes back as Type "a"
		t = askQuestion(qIndex)
	}

	if guessAnimal(t.Index) {
		fmt.Println("Why not try another animal?")
	} else {
		//Guess was wrong, let's learn a new animal
		learnNewAnimal(qIndex, t.Index)
	}
}

func askQuestion(i int) *Target {
	var t *Target

	q := questions[i]

	if askYesNoQuestion(q.Question + "?") {
		t = &q.YesTarget
	} else {
		t = &q.NoTarget
	}

	return t
}

func guessAnimal(i int) bool {
	return askYesNoQuestion("Is your animal " + withArticle(animals[i]) + "?")
}

func learnNewAnimal(qi, ai int) { //(this is kinda nasty but it works)
	fmt.Println("What is the animal you were thinking of called?") //nicer grammar than original
	a := getInput()

	//add to the animals array so it has an Index
	newAi := len(animals)
	animals = append(animals, a)

	//build a new question
	fmt.Println("Please type a Yes/No question that would distinguish " + withArticle(a) + " from " + withArticle(animals[ai]) + ":")
	q := getInput()

	//find out which route
	var yesTarget Target
	var noTarget Target

	if askYesNoQuestion("Please type the answer for " + withArticle(a) + ":") {
		yesTarget = Target{Type: ANSWER, Index: newAi}
		noTarget = Target{Type: ANSWER, Index: ai}
	} else {
		noTarget = Target{Type: ANSWER, Index: newAi}
		yesTarget = Target{Type: ANSWER, Index: ai}
	}

	q = strings.TrimRight(q, "?") //we add our own ? during question time

	//add the new question
	newQi := len(questions)
	questions = append(questions, Question{Question: q, YesTarget: yesTarget, NoTarget: noTarget})

	//amend the old question
	if questions[qi].YesTarget.Type == ANSWER && questions[qi].YesTarget.Index == ai {
		questions[qi].YesTarget.Type = QUESTION
		questions[qi].YesTarget.Index = newQi
	} else {
		questions[qi].NoTarget.Type = QUESTION
		questions[qi].NoTarget.Index = newQi
	}
}

func getArticle(noun string) string {
	if strings.ContainsAny(noun[:1], "AEIOUHaeiouh") {
		return "an"
	}

	return "a"
}

func withArticle(noun string) string {
	return getArticle(noun) + " " + noun
}

func getInput() (input string) {
	input, _ = reader.ReadString('\n')
	input = strings.TrimRight(input, "\r\n")
	return
}

func askYesNoQuestion(question string) bool {
	for true {
		fmt.Println(question)

		switch strings.ToLower(getInput())[:1] {
		case "n": //No
			return false
		case "y": //Yes
			return true
		default:
			fmt.Println("Please answer yes or no.")
		}
	}

	//this is impossible to reach, as the loop is infinite unless return is triggered by a case
	//but golint wants a return statement
	return false
}
