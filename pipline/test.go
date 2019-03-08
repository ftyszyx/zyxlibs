package pipline

import (
	"fmt"
)

func loadCheckpoint() {
	fmt.Println(fmt.Sprintf("loadCheckpoint"))
}

func extractReviewsFromA() {
	fmt.Println(fmt.Sprintf("loadCheckpoint"))
}

/*
func test() {
	//恢复上次执行的checkpoint，如果是第一次执行就获取一个初始值。
	checkpoint := loadCheckpoint()

	//工序(1)在pipeline外执行，最后一个工序是保存checkpoint
	pipeline := NewPipeline(4, 8, 2, 1)
	for {
		//(1)
		//加载100条数据，并修改变量checkpoint
		//data是数组，每个元素是一条评论，之后的联表、NLP都直接修改data里的每条记录。
		data, err := extractReviewsFromA(&checkpoint, 100)
		if err != nil {
			log.Print(err)
			break
		}

		//这里有个Golang著名的坑。
		//“checkpoint”是循环体外的变量，它在内存中只有一个实例并在循环中不断被修改，所以不能在异步中使用它。
		//这里创建一个副本curCheckpoint，储存本次循环的checkpoint。
		curCheckpoint := checkpoint

		ok := pipeline.Async(func() error {
			//(2)
			return joinUserFromB(data)
		}, func() error {
			//(3)
			return nlp(data)
		}, func() error {
			//(4)
			return loadDataToC(data)
		}, func() error {
			//(5)保存checkpoint
			log.Print("done:", curCheckpoint)
			return saveCheckpoint(curCheckpoint)
		})
		if !ok {
			break
		}

		if len(data) < 100 {
			break
		} //处理完毕
	}
	err := pipeline.Wait()
	if err != nil {
		fmt.Println(fmt.Sprintf("%v", err))
	}
}
*/
