package storage

import "math/rand"

type WOW struct {
	data []string
	size int
}

func NewWOW() WOW {
	data := []string{
		"Our Pious Predecessors – Living Examples of the Sunnah",
		"Advantages of adopting the appearance of the pious",
		"Studying the Lives of the Pious Predecessors Leads One to the Sunnah",
		"Valuing the Company of the Ahlullah",
		"Acquiring the True Love of Allah Ta‟ala in the Company of the Ahlullah",
		"The Preservation of One‟s Deen Depends on Remaining in the Company of the Ahlullah",
		"Obedience to Parents is a means of a Life of Contentmen",
		"Fulfilling the Rights of Parents after their Demise",
		"Making an Effort to Improve one‟s Character",
		"Despondency – The Trap of Shaitaan",
	}
	size := len(data)

	return WOW{
		data: data,
		size: size,
	}
}

func (w WOW) LoadPhrase() string {
	return w.data[rand.Intn(w.size)] //nolint:gosec
}
