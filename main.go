package main

import "fmt"

func main() {

	// nodes, err := GetNodes()
	// if err != nil {
	// 	fmt.Printf("Error retrieving nodes: %v", err)
	// 	fmt.Println()
	// 	return
	// }

	// pods, err := GetRunningPods()
	// if err != nil {
	// 	fmt.Printf("Error retrieving pods: %v", err)
	// 	fmt.Println()
	// 	return
	// }

	// unPods, err := GetUnschedulePods()
	// if err != nil {
	// 	fmt.Printf("Error retrieving unscheduled pods: %v", err)
	// 	fmt.Println()
	// 	return
	// }

	// for _, node := range nodes.Items {
	// 	fmt.Printf("Node Name: %s\n", node.Metadata.Name)
	// }

	// fmt.Println()

	// for _, pod := range pods.Items {
	// 	fmt.Printf("Running Pod Name: %s\n", pod.Metadata.Name)
	// }

	// fmt.Println()

	// for _, pod := range unPods {
	// 	fmt.Printf("Unscheduled Pod Name: %s\n", pod.Metadata.Name)
	// }

	fmt.Println("Starting mikischeduler...")

	pods, err := GetUnschedulePods()
	if err != nil {
		fmt.Println(err)
	}

	if pods == nil {
		fmt.Println("No pods to schedule.")
	}

	for _, pod := range pods {

		nodes, err := Fit(pod)
		if err != nil {
			fmt.Println(err)
		}

		for _, node := range nodes {
			fmt.Println("Node: " + node.Metadata.Name + " is reconsidering.")
		}

		node, err := BestPrice(nodes)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(node.Metadata.Name)

		err = Bind(pod, node)
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println("Pod: " + pod.Metadata.Name + " assignet to Node: " + node.Metadata.Name)

	}
}
