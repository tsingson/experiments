package main

func main() {

	d := diodes.NewManyToOne(1000, diodes.AlertFunc(func(missed int) {
		fmt.Printf("Dropped %d messages\n", missed)
	}))
	w := diode.NewWriter(os.Stdout, d, 10*time.Millisecond)
	log := zerolog.New(w)
	log.Print("test")
}
