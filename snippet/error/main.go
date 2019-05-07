package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
)

// 获取和返回已知的操作系统相关错误的潜在错误值
// underlyingError 会返回已知的操作系统相关错误的潜在错误值。
func underlyingError(err error) error {
	switch err := err.(type) {
	case *os.PathError:
		return err.Err
	case *os.LinkError:
		return err.Err
	case *os.SyscallError:
		return err.Err
	case *exec.Error:
		return err.Err
	}
	return err
}

// Errno 代表某种错误的类型。
type Errno int

func (e Errno) Error() string {
	return "errno " + strconv.Itoa(int(e))
}

func main() {
	// 错误处理方式

	//1. 对于类型在已知范围内的一系列错误值，一般使用类型断言表达式或类型switch语句来判断；
	//2. 对于已有相应变量且类型相同的一系列错误值，一般直接使用判等操作来判断；
	//3. 对于没有相应变量且类型未知的一系列错误值，只能使用其错误信息的字符串表示形式来做判断。
	{

		// 示例1。
		r, w, err := os.Pipe()
		if err != nil {
			fmt.Printf("unexpected error: %s\n", err)
			return
		}
		// 人为制造 *os.PathError 类型的错误。
		r.Close()
		_, err = w.Write([]byte("hi"))
		uError := underlyingError(err)
		fmt.Printf("underlying error: %s (type: %T)\n",
			uError, uError)
		fmt.Println()

		// 示例2。
		paths := []string{
			os.Args[0],           // 当前的源码文件或可执行文件。
			"/it/must/not/exist", // 肯定不存在的目录。
			os.DevNull,           // 肯定存在的目录。
		}
		printError := func(i int, err error) {
			if err == nil {
				fmt.Println("nil error")
				return
			}
			err = underlyingError(err)
			switch err {
			case os.ErrClosed:
				fmt.Printf("error(closed)[%d]: %s\n", i, err)
			case os.ErrInvalid:
				fmt.Printf("error(invalid)[%d]: %s\n", i, err)
			case os.ErrPermission:
				fmt.Printf("error(permission)[%d]: %s\n", i, err)
			}
		}
		var f *os.File
		var index int
		{
			index = 0
			f, err = os.Open(paths[index])
			if err != nil {
				fmt.Printf("unexpected error: %s\n", err)
				return
			}
			// 人为制造潜在错误为 os.ErrClosed 的错误。
			f.Close()
			_, err = f.Read([]byte{})
			printError(index, err)
		}
		{
			index = 1
			// 人为制造 os.ErrInvalid 错误。
			f, _ = os.Open(paths[index])
			_, err = f.Stat()
			printError(index, err)
		}
		{
			index = 2
			// 人为制造潜在错误为 os.ErrPermission 的错误。
			_, err = exec.LookPath(paths[index])
			printError(index, err)
		}
		if f != nil {
			f.Close()
		}
		fmt.Println()

		// 示例3。
		paths2 := []string{
			runtime.GOROOT(),     // 当前环境下的Go语言根目录。
			"/it/must/not/exist", // 肯定不存在的目录。
			os.DevNull,           // 肯定存在的目录。
		}
		printError2 := func(i int, err error) {
			if err == nil {
				fmt.Println("nil error")
				return
			}
			err = underlyingError(err)
			if os.IsExist(err) {
				fmt.Printf("error(exist)[%d]: %s\n", i, err)
			} else if os.IsNotExist(err) {
				fmt.Printf("error(not exist)[%d]: %s\n", i, err)
			} else if os.IsPermission(err) {
				fmt.Printf("error(permission)[%d]: %s\n", i, err)
			} else {
				fmt.Printf("error(other)[%d]: %s\n", i, err)
			}
		}
		{
			index = 0
			err = os.Mkdir(paths2[index], 0700)
			printError2(index, err)
		}
		{
			index = 1
			f, err = os.Open(paths[index])
			printError2(index, err)
		}
		{
			index = 2
			_, err = exec.LookPath(paths[index])
			printError2(index, err)
		}
		if f != nil {
			f.Close()
		}
	}

	// 错误处理指导方法
	// 1. 用类型建立起树形结构的错误体系，用统一字段建立起可追根溯源的链式错误关联
	// 2. 扁平化错误号
	// 潜在问题：错误值可能被修改
	// 解决方法：
	// 	2.1 私有化错误号，提供方法来判断错误值
	//  2.2 自定义错误类型，uintptr
	{
		var err error
		// 示例1。
		_, err = exec.LookPath(os.DevNull)
		fmt.Printf("error: %s\n", err)
		if execErr, ok := err.(*exec.Error); ok {
			execErr.Name = os.TempDir()
			execErr.Err = os.ErrNotExist
		}
		fmt.Printf("error: %s\n", err)
		fmt.Println()

		// 示例2。
		err = os.ErrPermission
		if os.IsPermission(err) {
			fmt.Printf("error(permission): %s\n", err)
		} else {
			fmt.Printf("error(other): %s\n", err)
		}
		os.ErrPermission = os.ErrExist
		// 上面这行代码修改了os包中已定义的错误值。
		// 这样做会导致下面判断的结果不正确。
		// 并且，这会影响到当前Go程序中所有的此类判断。
		// 所以，一定要避免这样做！
		if os.IsPermission(err) {
			fmt.Printf("error(permission): %s\n", err)
		} else {
			fmt.Printf("error(other): %s\n", err)
		}
		fmt.Println()

		// 示例3。
		const (
			ERR0 = Errno(0)
			ERR1 = Errno(1)
			ERR2 = Errno(2)
		)
		var myErr error = Errno(0)
		switch myErr {
		case ERR0:
			fmt.Println("ERR0")
		case ERR1:
			fmt.Println("ERR1")
		case ERR2:
			fmt.Println("ERR2")
		}
	}
}
