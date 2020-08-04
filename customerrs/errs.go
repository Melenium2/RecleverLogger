package customerrs

import "fmt"

func ServerConfigIsNilErr() error {
	return fmt.Errorf(`[MEW-1] Server config is nil`)
}

func ServerHttpPortIsEmptyErr() error {
	return fmt.Errorf(`[MEW-2] Server http port is empty`)
}

func ServerFailToListenPortErr(port string, err error) error {
	return fmt.Errorf(`[MEW-3] Server failed to create new TCP listener port = %s, err = %v`, port, err)
}

func ServerHaveNoHandlersErr() error {
	return fmt.Errorf(`[MEW-4] Server have no handlers`)
}

func ServerFailedToShutdownErr() error {
	return fmt.Errorf(`[MEW-5] Server failed to shutdown`)
}

func ServiceConfigIsNilErr() error {
	return fmt.Errorf(`[MEW-6] Service config is nil`)
}

func ServiceConfigDbUrlIsEmptyErr() error {
	return fmt.Errorf(`[MEW-6] Service config db url is empty`)
}