"use client"

import { useState, useEffect } from "react"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Input } from "@/components/ui/input"
import { Textarea } from "@/components/ui/textarea"
import { Badge } from "@/components/ui/badge"
import { Checkbox } from "@/components/ui/checkbox"
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select"
import { Plus, Calendar, User } from "lucide-react"
import { api, Task } from "@/lib/api"


interface DealTasksProps {
  dealId: number
}

export function DealTasks({ dealId }: DealTasksProps) {
  const [tasks, setTasks] = useState<Task[]>([])
  const [isLoading, setIsLoading] = useState(true)
  const [isCreating, setIsCreating] = useState(false)
  const [formData, setFormData] = useState({
    task_name: "",
    task_description: "",
    due_date: "",
    assigned_to: "",
  })

  useEffect(() => {
    loadTasks()
  }, [dealId])

  const loadTasks = async () => {
    try {
      setIsLoading(true)
      const tasksData = await api.getTasksForDeal(dealId)
      setTasks(tasksData)
    } catch (err) {
      console.error("Failed to load tasks:", err)
    } finally {
      setIsLoading(false)
    }
  }

  const handleCreateTask = async () => {
    if (!formData.task_name.trim()) return

    try {
      await api.createTaskForDeal(dealId, {
        task_name: formData.task_name,
        task_description: formData.task_description,
        due_date: formData.due_date,
        assigned_to: Number(formData.assigned_to),
        status: "Pending",
      })
      setFormData({ task_name: "", task_description: "", due_date: "", assigned_to: "" })
      setIsCreating(false)
      loadTasks()
    } catch (err) {
      console.error("Failed to create task:", err)
    }
  }

  const handleToggleTask = async (taskId: number, currentStatus: string) => {
    const newStatus = currentStatus === "Pending" ? "Completed" : "Pending"

    try {
      await api.updateTaskForDeal(dealId, taskId, { status: newStatus })
      loadTasks()
    } catch (err) {
      console.error("Failed to update task:", err)
    }
  }

  const isOverdue = (dueDate: string) => {
    return new Date(dueDate) < new Date() && new Date(dueDate).toDateString() !== new Date().toDateString()
  }

  const getStatusBadge = (task: Task) => {
    if (task.status === "Completed") {
      return (
        <Badge variant="outline" className="bg-green-50 text-green-700 border-green-200">
          Completed
        </Badge>
      )
    }
    if (isOverdue(task.due_date)) {
      return (
        <Badge variant="outline" className="bg-red-50 text-red-700 border-red-200">
          Overdue
        </Badge>
      )
    }
    return (
      <Badge variant="outline" className="bg-yellow-50 text-yellow-700 border-yellow-200">
        Pending
      </Badge>
    )
  }

  if (isLoading) {
    return <div className="text-center py-8">Loading tasks...</div>
  }

  return (
    <div className="space-y-4">
      {/* Create Task */}
      <Card>
        <CardHeader>
          <div className="flex items-center justify-between">
            <div>
              <CardTitle>Deal Tasks</CardTitle>
              <CardDescription>Manage tasks and action items for this deal</CardDescription>
            </div>
            {!isCreating && (
              <Button onClick={() => setIsCreating(true)} size="sm" className="bg-cyan-600 hover:bg-cyan-700">
                <Plus className="mr-2 h-4 w-4" />
                Add Task
              </Button>
            )}
          </div>
        </CardHeader>
        {isCreating && (
          <CardContent className="space-y-4">
            <div className="grid gap-4 md:grid-cols-2">
              <div>
                <Input
                  placeholder="Task name"
                  value={formData.task_name}
                  onChange={(e) => setFormData((prev) => ({ ...prev, task_name: e.target.value }))}
                />
              </div>
              <div>
                <Input
                  type="date"
                  value={formData.due_date}
                  onChange={(e) => setFormData((prev) => ({ ...prev, due_date: e.target.value }))}
                />
              </div>
            </div>
            <Textarea
              placeholder="Task description (optional)"
              value={formData.task_description}
              onChange={(e) => setFormData((prev) => ({ ...prev, task_description: e.target.value }))}
              rows={2}
            />
            <Select
              value={formData.assigned_to}
              onValueChange={(value) => setFormData((prev) => ({ ...prev, assigned_to: value }))}
            >
              <SelectTrigger>
                <SelectValue placeholder="Assign to user" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="1">User 1</SelectItem>
                <SelectItem value="2">User 2</SelectItem>
                <SelectItem value="3">User 3</SelectItem>
              </SelectContent>
            </Select>
            <div className="flex gap-2">
              <Button onClick={handleCreateTask} size="sm" className="bg-cyan-600 hover:bg-cyan-700">
                Create Task
              </Button>
              <Button
                onClick={() => {
                  setIsCreating(false)
                  setFormData({ task_name: "", task_description: "", due_date: "", assigned_to: "" })
                }}
                variant="outline"
                size="sm"
              >
                Cancel
              </Button>
            </div>
          </CardContent>
        )}
      </Card>

      {/* Tasks List */}
      <div className="space-y-4">
        {tasks.length === 0 ? (
          <Card>
            <CardContent className="text-center py-8">
              <p className="text-muted-foreground">No tasks yet. Add your first task to get started.</p>
            </CardContent>
          </Card>
        ) : (
          tasks.map((task) => (
            <Card key={task.id} className={task.status === "Completed" ? "opacity-75" : ""}>
              <CardContent className="pt-6">
                <div className="flex items-start gap-4">
                  <Checkbox
                    checked={task.status === "Completed"}
                    onCheckedChange={() => handleToggleTask(task.id, task.status)}
                    className="mt-1"
                  />
                  <div className="flex-1 space-y-2">
                    <div className="flex items-center justify-between">
                      <h4
                        className={`font-medium ${task.status === "Completed" ? "line-through text-muted-foreground" : ""}`}
                      >
                        {task.task_name}
                      </h4>
                      {getStatusBadge(task)}
                    </div>
                    {task.task_description && <p className="text-sm text-muted-foreground">{task.task_description}</p>}
                    <div className="flex items-center gap-4 text-sm text-muted-foreground">
                      <div className="flex items-center gap-1">
                        <Calendar className="h-4 w-4" />
                        Due: {new Date(task.due_date).toLocaleDateString()}
                      </div>
                      <div className="flex items-center gap-1">
                        <User className="h-4 w-4" />
                        Assigned to User {task.assigned_to}
                      </div>
                    </div>
                  </div>
                </div>
              </CardContent>
            </Card>
          ))
        )}
      </div>
    </div>
  )
}
