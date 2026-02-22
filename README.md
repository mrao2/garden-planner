# garden-planner

Docs:
To create a prioritized queue that ranks garden vegetables based on when they need to be planted outside in USDA zone 6b, you can use a min heap with each vegetable being represented as a node in the heap.

The priority of each vegetable will be determined by the number of days before the last frost date in USDA zone 6b that it needs to be planted. The vegetable with the smallest number of days before the last frost date will have the highest priority and will be placed at the root of the heap.

Here are the steps to create the prioritized queue:

    1. Define a node class to represent the vegetables in the queue. Each node should have two attributes: the name of the vegetable and the number of days before the last frost date that it needs to be planted.

    2. Create an empty min heap to represent the queue.

    3. For each vegetable, create a node and add it to the heap using the number of days before the last frost date as the priority.

    4. To extract the highest priority vegetable from the queue, remove the root node of the heap.

    5. To add a new vegetable to the queue, create a node for the vegetable and insert it into the heap using the number of days before the last frost date as the priority.

    6. To maintain the heap property after inserting or removing a node, use the heapify function to restore the heap property.

To sort the priority of vegetables based on additional factors such as maturity time, successive sowings, and sunlight requirements, you can modify the Vegetable class in the previous implementation to include attributes for these factors. Then, you can modify the comparison function of the class to take into account all the factors when comparing two vegetables.