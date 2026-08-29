package config

const (
	UsernameAdjectivesKey = "username_adjectives"
	UsernameAnimalsKey    = "username_animals"
	UsernameSpaceWordKey  = "username_space_word"
)

var defaultUsernameAdjectives = []string{
	"swift", "quiet", "bold", "crisp", "bright", "clever", "calm", "eager", "brave", "nimble",
	"rapid", "gentle", "merry", "proud", "keen", "lively", "sunny", "lucky", "daring", "noble",
	"prime", "vivid", "sleek", "grand", "wise", "quick", "steady", "jolly", "brisk", "fresh",
}

var defaultUsernameAnimals = []string{
	"ant", "ape", "bat", "bear", "bee", "bison", "cat", "clam", "cobra", "crab",
	"crow", "deer", "dove", "duck", "eagle", "eel", "elk", "falcon", "finch", "fox",
	"frog", "gecko", "goat", "goose", "hare", "hawk", "heron", "horse", "ibis", "jay",
	"koala", "lemur", "lion", "llama", "lynx", "mole", "mouse", "otter", "owl", "panda",
	"quail", "raven", "seal", "shark", "sheep", "skunk", "sloth", "tiger", "wolf", "zebra",
}

var defaultUsernameSpaceWords = []string{
	"star", "comet", "nova", "orbit", "lunar", "solar", "meteor", "nebula", "galaxy", "planet",
	"cosmos", "quasar", "pulsar", "rocket", "saturn", "mars", "venus", "pluto", "eclipse", "flare",
	"zenith", "aurora", "crater", "launch", "module", "probe", "lander", "shuttle", "apollo", "orion",
}

func DefaultUsernameAdjectives() []string {
	return append([]string(nil), defaultUsernameAdjectives...)
}

func DefaultUsernameAnimals() []string {
	return append([]string(nil), defaultUsernameAnimals...)
}

func DefaultUsernameSpaceWords() []string {
	return append([]string(nil), defaultUsernameSpaceWords...)
}
