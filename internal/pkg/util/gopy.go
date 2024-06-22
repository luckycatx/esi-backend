package util

// func ExecPython(cmd string) error {
// 	// e.g. cmd = "./internal/util/pytest.py"
// 	py := exec.Command("python", cmd)
// 	output, err := py.CombinedOutput()
// 	if err != nil {
// 		fmt.Printf("Error when executing python: %s\n", err)
// 		return fmt.Errorf("error executing python: %w", err)
// 	}
// 	fmt.Print(string(output))
// 	return nil
// }

// func GoPy(c *gin.Context) {
// 	if err := ExecPython("./internal/util/pytest.py"); err != nil {
// 		dosomething
// 	}
// }
