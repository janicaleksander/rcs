package createtaskstate

func (c *CreateTaskScene) Reset() {
	c.errorSection.isSetupError = false
	c.errorSection.isCreateError = false
	c.infoSection.isInfoMessage = false
	c.infoSection.infoMessage = ""
	c.errorSection.errorMessage = ""
}
func (c *CreateTaskScene) FetchUsers() {
	//TODO get proper lvl value
	/*
		res, err := utils.MakeRequest(utils.NewRequest(c.cfg.Ctx, c.cfg.ServerPID, &proto.GetUserAboveLVL{
			Lower: -1,
			Upper: 10,
		})) //TODO
		if err != nil {
			//context deadline exceeded
			//do sth with that
			c.errorSection.isSetupError = true
		}

		if v, ok := res.(*proto.UsersAboveLVL); ok {
			c.newUnitSection.usersDropdown.Strings = make([]string, 0, 64)
			c.newUnitSection.usersDropdown.Strings = append(c.newUnitSection.usersDropdown.Strings,
				"Choose user by his ID")
			for _, user := range v.Users {
				c.newUnitSection.usersDropdown.Strings = append(c.newUnitSection.usersDropdown.Strings,
					user.Id+"\n"+user.Email)
			}
		} else {
			c.errorSection.isSetupError = true
		}
	*/
}
