package pkg

type RoundCounter uint

type Gofeur struct {
	Ui *UI
	rc RoundCounter
}

func (gofeur *Gofeur) Init() {
	gofeur.Ui = UIStart()
}

func (gofeur *Gofeur) Run() {
	feurUI := gofeur.Ui

	if err := feurUI.App.SetRoot(feurUI.Layout, true).
                         SetFocus(feurUI.OutputBox).
                         Run(); err != nil {
		panic(err)
	}
}
