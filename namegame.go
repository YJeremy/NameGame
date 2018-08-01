// This sample program demonstrates how to use an unbuffered
// channel to simulate a game of tennis between two goroutines.
package main

import (
	"fmt"
	"github.com/urfave/cli"
	"math/rand"
	"os"
	"sync"
	"time"
)

//创建计数器对象
var wg sync.WaitGroup

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	printExplain() //打印一些游戏规则

	court := make(chan int)
	app := cli.NewApp()
	app.Name = "Name battle"
	app.Usage = "Using two name play game"
	app.Version = "0.0"

	app.Flags = []cli.Flag{ //创建应用的参数cli.Flag类型
		cli.StringFlag{
			Name:  "operation,o",
			Value: "h",
			Usage: "opertaion",
		},
		cli.StringFlag{
			Name:  "name1,n1",
			Value: "",
			Usage: "user name1",
		},
		cli.StringFlag{
			Name:  "name2,n2",
			Value: "",
			Usage: "user name2",
		},
		cli.Int64Flag{
			Name:  "name1HP,hp1",
			Value: 20,
			Usage: "user name1 HP",
		},
		cli.Int64Flag{
			Name:  "name2HP,hp2",
			Value: 20,
			Usage: "user name2 HP",
		},
	}

	app.Action = func(c *cli.Context) error {

		operation := c.String("operation")
		name1 := c.String("name1")
		name2 := c.String("name2")
		name1hp := c.Int("name1HP")
		name2hp := c.Int("name2HP")
		fmt.Printf("游戏开始%s血量值为%d\t%s血量值为%d\n", name1, name1hp, name2, name2hp)
		// 增加并发计数
		wg.Add(2)

		if operation == "p" {
			go player(name1, court, 1, name1hp)
			go player(name2, court, 1, name2hp)

			// 给通道传入数据，触发并发开始游戏
			court <- 1

			// 等候计数器完结
			wg.Wait()
		}
		return nil
	}
	app.Run(os.Args)

}

// 游戏主要逻辑
func player(name string, court chan int, turn, hp int) {
	// 函数结束时，关闭计数器；闭包传值
	defer wg.Done()
	skill := []string{"会心一击", "阳光烈焰", "天马流星拳", "升龙拳", "吸血光环", "十万伏特", "射火焰", "还我漂漂拳", "十字死光", "百万吨吸收", "挠痒痒 ", "狮子偷桃", "夺命剪刀脚", "龟波气功"}
	hecai := []string{"厉害厉害～", "优秀！", "skr～", "啊哟，不错哦！", "一个字，稳！"}

	for {

		battle, ok := <-court

		if !ok {

			fmt.Printf("%s输了,555～\n", name)
			return
		}

		n := rand.Intn(10)
		hecainum := rand.Intn(4)

		if hp < 0 {
			hp = 0
			fmt.Printf("KO！%s将对方血量打到0！ %s胜利！\n", name, name)

			// 关闭通道，退出循环
			close(court)
			return
		}

		fmt.Printf("第%d回合 ", turn)
		turn++
		if n == 4 {
			hp = -1
			fmt.Printf("你的名字|%-6s|使出秘籍“要你命3000”！直接将对手血量打成0！\n%s胜利！\n", name, name)
			close(court)
			return
		} else {

			hp = hp - n
			fmt.Printf("你的名字|%-6s|使用|%s输出%d点伤害,%-6s|\t对手血量%d\n", name, skill[n], n, hecai[hecainum], hp)
		}
		time.Sleep(1 * time.Second)
		court <- battle //阻塞通道
	}
}
func printExplain() {
	fmt.Println("输入 -h 获取参数帮助")
	fmt.Println("===================")
	fmt.Println("快速游戏：-n1 名字1 -n2 名字2")
	fmt.Println("开始默认血量20,可以通过参数 -hp1 -hp2 修改两人血量")
	fmt.Println("\n")
}
