# gotdots-cli
Easy to use application to share dotfiles as packages!
----
Create packages to share your desktop rices.

Help: dots [command] [options]
	
	new      <pack name> :	Creates new package with given name
	update	 <pack name> :  Creates new version for existing package
	push     <pack name> :	Pushes package to registry (aws s3 for now)
	get      <pack name> :	Downloads package
	install  <pack name> :	Installs package to system
