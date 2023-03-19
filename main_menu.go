package main

type MainMenu struct {
	GenericWidgetContainerScreen
}

func (s *MainMenu) LoadResources() {
	rm := ResourceManager_GetInstance()

	rm.LoadImage("assets/gui/button.png")
	rm.LoadImage("assets/gui/gui_frame_herb.png")
	rm.LoadSound("sound/bgm_main_menu.ogg")
}

func (s *MainMenu) ProcessKeyEvents() bool {
	return s.GenericWidgetContainerScreen.ProcessKeyEvents()
}

func (s *MainMenu) Update() {
	s.GenericWidgetContainerScreen.Update()
	s.ProcessKeyEvents()
}

func CreateMainMenu(g *Game) *MainMenu {
	s := new(MainMenu)
	s.IScreen = s
	s.game = g
	s.title = "Main Menu"

	s.widgets = append(s.widgets, CreateCommonButton(g, I18n("string_new_game", "New Game"), func(*Game) {
		g.SetScreen(NewGameplayScreen(g))
	}))

	s.widgets = append(s.widgets, CreateCommonButton(g, I18n("string_load_save", "Load Save"), func(*Game) {

	}))

	s.widgets = append(s.widgets, CreateCommonButton(g, I18n("string_settings", "Settings"), func(*Game) {
		//g.SetScreen(NewSettingsScreen(g))
	}))

	s.widgets = append(s.widgets, CreateCommonButton(g, I18n("string_exit", "Exit"), func(*Game) {
		g.SetScreen(NewFarewellScreen(g))
	}))

	s.SetInitialFocus()

	return s
}
