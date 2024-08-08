package metamigrator

// MetaMigratorTask is the task of identifying migration data and then migrating that data to cloud.
type MetaMigratorTask struct {
	service *MetaMigratorService
}

// NewMetaMigratorTask returns a new MetaMigratorTask instance.
func NewMetaMigratorTask(service *MetaMigratorService) *MetaMigratorTask {
	return &MetaMigratorTask{
		service: service,
	}
}
