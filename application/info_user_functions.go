package application

func (i *InfoUserScene) Reset() {
	i.lastProcessedUserIdx = -1
	i.descriptionName = ""
	i.descriptionSurname = ""
	i.descriptionLVL = ""
	i.currSelectedUserID = ""

}
