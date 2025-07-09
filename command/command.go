package command

import (
	"log"
	"os"
	"os/exec"
)

// ExecuteShellCommandAsync 异步执行 shell 命令。
// 该函数会发送一个信号到 maxShellExecConcurrent 通道，以控制并发执行的最大数量。
// 参数:
//
//	command: 要执行的 shell 命令。
//	maxShellExecConcurrent: 用于控制最大并发执行数的通道。
func ExecuteShellCommandAsync(command string, maxShellExecConcurrent chan struct{}) {
	// 发送一个信号到通道，表示开始执行一个命令。
	maxShellExecConcurrent <- struct{}{}

	go func() {
		// 执行 shell 命令
		cmd := exec.Command("sh", "-c", command)
		// 不能缺少 HOME=/root 否则会报错
		cmd.Env = append(os.Environ(), "HOME=/root")
		// 获取输出结果
		output, err := cmd.CombinedOutput()
		if err != nil {
			log.Printf("执行shell命令失败，失败原因: %s", output)
			return
		}
		log.Printf("执行成功，输出结果: %s", output)

		// 接收一个信号，表示命令执行完成，释放并发控制资源。
		<-maxShellExecConcurrent
	}()
}
