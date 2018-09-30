package main

import (
	"math/rand"
	"time"
)

type Joke struct{}

var jokes = []string{
	"What time did the man go to the dentist? Tooth hurt-y.",
	"Did you hear about the guy who invented Lifesavers? They say he made a mint.",
	"A ham sandwich walks into a bar and orders a beer. Bartender says, 'Sorry we don't serve food here.'",
	"Why do chicken coops only have two doors? Because if they had four, they would be chicken sedans!",
	"Why did the Clydesdale give the pony a glass of water? Because he was a little horse!",
	"How do you make a Kleenex dance? Put a little boogie in it!",
	"Why do they put fences around graveyards? Everyone is dying to get in.",
	"Two peanuts were walking down the street. One was a salted.",
	"I used to have a job at a calendar factory, but I got fired when I took a couple of days off.",
	"How do you make holy water? You boil the hell out of it.",
	"Two guys walk into a bar, the third one ducks.",
	"I had a dream that I was a muffler last night. I woke up exhausted!",
	"Our wedding was so beautiful even the cake was in tiers.",
	"I'm reading a book on glue I just can't seem to put it down.",
	"What do you call an Argentinian with a rubber toe? Roberto.",
	"I am terrified of elevators I'm going to start taking steps to avoid them.",
	"Why do crabs never give to charity? Because they're shellfish..",
	"Why don't skeletons ever go trick or treating? Because they have no body to go with.",
	"What do you call cheese by itself? Provolone.",
	"'Ill call you later.' Don't call me later, call me Dad.",
	"Dad, I'm hungry. Hello, Hungry. I'm Dad.",
	"Where does Fonzie like to go for lunch? Chick-Fil-Eyyyyyyyy.",
	"Did you hear about the cheese factory that exploded in France? There was nothing left but de Brie.",
	"I knew I shouldn't have had the seafood I'm feeling a little eel.",
	"What do you call a sketchy Italian neighbour hood? The Spaghetto.",
	"Why can't you have a nose 12 inches long? Because then it would be a foot.",
	"My wife is on a tropical food diet, the house is full of the stuff It's enough to make a mango crazy.",
	"I'd like to give a big shout out to all the sidewalks for keeping me off the streets.",
	"What does an annoying pepper do? It get's jalapeno face.",
	"Why did the scarecrow win an award? Because he was outstanding in his field.",
	"Why do bees hum? Because they don't know the words.",
	"What do prisoners use to call each other? Cell phones.",
	"What do you call cheese that isn't yours? Nacho Cheese.",
	"What do you get when you cross a snowman with a vampire? Frostbite.",
	"What lies at the bottom of the ocean and twitches? A nervous wreck.",
	"Why couldn't the bicycle stand up by itself? It was two tired.",
	"What did the grape do when he got stepped on? He let out a little wine.",
	"I've just been diagnosed as colorblind I know, it certainly has come out of the purple.",
	"Last night I dreamt I was a muffler I woke up exhausted.",
	"What do you call a deer with no eyes? No idea.",
	"I just watched a program about beavers It was the best dam program I've ever seen.",
	"Two satellites decided to get married The wedding wasn't much, but the reception was incredible.",
	"Did you hear about the guy who invented the knock knock joke? He won the 'no-bell' prize.",
	"Is this pool safe for diving? It deep ends.",
	"I used to hate facial hair but then it grew on me.",
	"What do you call a fake noodle? An Impasta.",
	"Can February March? No, but April May.",
	"Wanna hear a joke about paper? Nevermind, it's tearable.",
	"Don't trust atoms They make up everything.",
	"How many apples grow on a tree? All of them.",
	"What do you call somebody with no body and no nose? Nobody knows.",
	"What did the buffalo say to his son when he left for college? Bison.",
	"What do you call a pony with a sore throat? A little horse.",
	"I bought shoes from a drug dealer once I don't know what he laced them with, but I was tripping all day.",
	"Where do you learn to make ice cream? Sunday school.",
	"What did the officer molecule say to the suspect molecule? I've got my ion you.",
	"If prisoners could take their own mug shots they'd be called cellfies.",
	"Why can't you hear a pterodactyl go to the bathroom? Because the pee is silent.",
	"I'm not addicted to brake fluid I can stop whenever I want.",
	"Why did the coffee file a police report? It got mugged.",
	"Did you hear about the restaurant on the moon? Great food, no atmosphere.",
	"I hate jokes about german sausages They're the wurst.",
	"Why did the can-crusher quit his job? Because it was soda-pressing.",
	"I wouldn't buy anything with velcro It's a total rip-off.",
	"Dad, can you put the cat out? I didn't know it was on fire.",
	"This graveyard looks overcrowded people must be dying to get in there.",
	"Dad, can you put my shoes on? I don't think they'll fit me.",
	"Dad, did you get a haircut? No I got them all cut.",
	"Have you heard of the band 1023MB? They haven't got a gig yet.",
}

func (j Joke) Get() string {
	rand.Seed(time.Now().Unix())
	n := rand.Int() % len(jokes)
	return jokes[n]
}
