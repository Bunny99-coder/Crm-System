"use client"

import { useState, useEffect } from "react"
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Badge } from "@/components/ui/badge"
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog"
import { Label } from "@/components/ui/label"
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select"
import { Textarea } from "@/components/ui/textarea"
import { DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuTrigger } from "@/components/ui/dropdown-menu"
import { Alert, AlertDescription } from "@/components/ui/alert"
import { MoreHorizontal, Search, Plus, Edit, Trash2, CheckSquare, Clock, Calendar } from "lucide-react"
import { api, type Task, type User as ApiUser } from "@/lib/api"

export function TasksManagement() {
  const [tasks, setTasks] = useState<Task[]>([])
  const [users, setUsers] = useState<ApiUser[]>([])
  const [searchTerm, setSearchTerm] = useState("")
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState("")
  const [isCreateDialogOpen, setIsCreateDialogOpen] = useState(false)
  const [isEditDialogOpen, setIsEditDialogOpen] = useState(false)
  const [selectedTask, setSelectedTask] = useState<Task | null>(null)
  const [formData, setFormData] = useState({
    task_name: "",
    task_description: "",
    due_date: "",
    status: "Pending" as "Pending" | "Completed",
    assigned_to: "",
  })

  // Load data on component mount
  useEffect(() => {
    loadAllData()
  }, [])

  const loadAllData = async () => {
    try {
      setIsLoading(true)
      const [tasksData, usersData] = await Promise.all([api.getTasks(), api.getUsers()])
      setTasks(tasksData)
      setUsers(usersData)
    } catch (err) {
      setError("Failed to load data")
    } finally {
      setIsLoading(false)
    }
  }

  const handleCreateTask = async () => {
    try {
      const taskData = {
        task_name: formData.task_name,
        task_description: formData.task_description,
        due_date: formData.due_date,
        status: formData.status,
        assigned_to: Number.parseInt(formData.assigned_to),
      }
      await api.createTask(taskData)
      setIsCreateDialogOpen(false)
      resetForm()
      loadAllData()
    } catch (err) {
      setError("Failed to create task")
    }
  }

  const handleEditTask = async () => {
    if (!selectedTask?.id) return

    try {
      const taskData = {
        task_name: formData.task_name,
        task_description: formData.task_description,
        due_date: formData.due_date,
        status: formData.status,
        assigned_to: Number.parseInt(formData.assigned_to),
      }
      await api.updateTask(selectedTask.id, taskData)
      setIsEditDialogOpen(false)
      resetForm()
      setSelectedTask(null)
      loadAllData()
    } catch (err) {
      setError("Failed to update task")
    }
  }

  const handleDeleteTask = async (id: number) => {
    if (!confirm("Are you sure you want to delete this task?")) return

    try {
      await api.deleteTask(id)
      loadAllData()
    } catch (err) {
      setError("Failed to delete task")
    }
  }

  const handleToggleTaskStatus = async (task: Task) => {
    if (!task.id) return

    try {
      const newStatus = task.status === "Pending" ? "Completed" : "Pending"
      const taskData = {
        task_name: task.task_name,
        task_description: task.task_description,
        due_date: task.due_date,
        status: newStatus,
        assigned_to: task.assigned_to,
      }
      await api.updateTask(task.id, taskData)
      loadAllData()
    } catch (err) {
      setError("Failed to update task status")
    }
  }

  const resetForm = () => {
    setFormData({
      task_name: "",
      task_description: "",
      due_date: "",
      status: "Pending",
      assigned_to: "",
    })
  }

  const openEditDialog = (task: Task) => {
    setSelectedTask(task)
    setFormData({
      task_name: task.task_name,
      task_description: task.task_description,
      due_date: task.due_date.split("T")[0], // Format for date input
      status: task.status,
      assigned_to: task.assigned_to.toString(),
    })
    setIsEditDialogOpen(true)
  }

  const getUserName = (userId: number) => {
    const user = users.find((u) => u.id === userId)
    return user ? user.username : `User ${userId}`
  }

  const getStatusBadge = (status: string) => {
    switch (status) {
      case "Pending":
        return (
          <Badge variant="outline" className="bg-yellow-50 text-yellow-700 border-yellow-200">
            <Clock className="w-3 h-3 mr-1" />
            Pending
          </Badge>
        )
      case "Completed":
        return (
          <Badge variant="outline" className="bg-green-50 text-green-700 border-green-200">
            <CheckSquare className="w-3 h-3 mr-1" />
            Completed
          </Badge>
        )
      default:
        return <Badge variant="outline">{status}</Badge>
    }
  }

  const formatDate = (dateString: string) => {
    const date = new Date(dateString)
    return date.toLocaleDateString()
  }

  const isOverdue = (dueDate: string, status: string) => {
    if (status === "Completed") return false
    const today = new Date()
    const due = new Date(dueDate)
    return due < today
  }

  const filteredTasks = tasks.filter(
    (task) =>
      task.task_name.toLowerCase().includes(searchTerm.toLowerCase()) ||
      task.task_description.toLowerCase().includes(searchTerm.toLowerCase()) ||
      getUserName(task.assigned_to).toLowerCase().includes(searchTerm.toLowerCase()),
  )

  const pendingTasks = tasks.filter((t) => t.status === "Pending").length
  const completedTasks = tasks.filter((t) => t.status === "Completed").length
  const overdueTasks = tasks.filter((t) => isOverdue(t.due_date, t.status)).length

  if (isLoading) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="text-muted-foreground">Loading tasks...</div>
      </div>
    )
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div className="space-y-1">
          <h1 className="text-3xl font-bold text-balance text-foreground">Tasks Management</h1>
          <p className="text-muted-foreground text-pretty">Organize and track your team's tasks and deadlines</p>
        </div>
        <Dialog open={isCreateDialogOpen} onOpenChange={setIsCreateDialogOpen}>
          <DialogTrigger asChild>
            <Button className="bg-cyan-600 hover:bg-cyan-700">
              <Plus className="mr-2 h-4 w-4" />
              Add Task
            </Button>
          </DialogTrigger>
          <DialogContent className="sm:max-w-[500px]">
            <DialogHeader>
              <DialogTitle>Create New Task</DialogTitle>
              <DialogDescription>Add a new task to your team's workflow.</DialogDescription>
            </DialogHeader>
            <div className="grid gap-4 py-4">
              <div className="grid grid-cols-4 items-center gap-4">
                <Label htmlFor="task_name" className="text-right">
                  Task Name
                </Label>
                <Input
                  id="task_name"
                  value={formData.task_name}
                  onChange={(e) => setFormData((prev) => ({ ...prev, task_name: e.target.value }))}
                  className="col-span-3"
                  placeholder="Enter task name"
                  required
                />
              </div>
              <div className="grid grid-cols-4 items-center gap-4">
                <Label htmlFor="task_description" className="text-right">
                  Description
                </Label>
                <Textarea
                  id="task_description"
                  value={formData.task_description}
                  onChange={(e) => setFormData((prev) => ({ ...prev, task_description: e.target.value }))}
                  className="col-span-3"
                  placeholder="Task description"
                  rows={3}
                />
              </div>
              <div className="grid grid-cols-4 items-center gap-4">
                <Label htmlFor="due_date" className="text-right">
                  Due Date
                </Label>
                <Input
                  id="due_date"
                  type="datetime-local"
                  value={formData.due_date}
                  onChange={(e) => setFormData((prev) => ({ ...prev, due_date: e.target.value }))}
                  className="col-span-3"
                  required
                />
              </div>
              <div className="grid grid-cols-4 items-center gap-4">
                <Label htmlFor="assigned_to" className="text-right">
                  Assigned To
                </Label>
                <Select
                  value={formData.assigned_to}
                  onValueChange={(value) => setFormData((prev) => ({ ...prev, assigned_to: value }))}
                >
                  <SelectTrigger className="col-span-3">
                    <SelectValue placeholder="Select user" />
                  </SelectTrigger>
                  <SelectContent>
                    {users.map((user) => (
                      <SelectItem key={user.id} value={user.id!.toString()}>
                        {user.username}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </div>
              <div className="grid grid-cols-4 items-center gap-4">
                <Label htmlFor="status" className="text-right">
                  Status
                </Label>
                <Select
                  value={formData.status}
                  onValueChange={(value: "Pending" | "Completed") =>
                    setFormData((prev) => ({ ...prev, status: value }))
                  }
                >
                  <SelectTrigger className="col-span-3">
                    <SelectValue placeholder="Select status" />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="Pending">Pending</SelectItem>
                    <SelectItem value="Completed">Completed</SelectItem>
                  </SelectContent>
                </Select>
              </div>
            </div>
            <DialogFooter>
              <Button type="submit" onClick={handleCreateTask} className="bg-cyan-600 hover:bg-cyan-700">
                Create Task
              </Button>
            </DialogFooter>
          </DialogContent>
        </Dialog>
      </div>

      {error && (
        <Alert variant="destructive">
          <AlertDescription>{error}</AlertDescription>
        </Alert>
      )}

      {/* Stats */}
      <div className="grid gap-4 md:grid-cols-4">
        <div className="rounded-lg border bg-card p-6">
          <div className="flex items-center gap-2">
            <CheckSquare className="h-5 w-5 text-cyan-600" />
            <h3 className="font-semibold">Total Tasks</h3>
          </div>
          <p className="text-2xl font-bold text-cyan-600 mt-2">{tasks.length}</p>
        </div>
        <div className="rounded-lg border bg-card p-6">
          <div className="flex items-center gap-2">
            <Clock className="h-5 w-5 text-yellow-600" />
            <h3 className="font-semibold">Pending</h3>
          </div>
          <p className="text-2xl font-bold text-yellow-600 mt-2">{pendingTasks}</p>
        </div>
        <div className="rounded-lg border bg-card p-6">
          <div className="flex items-center gap-2">
            <CheckSquare className="h-5 w-5 text-green-600" />
            <h3 className="font-semibold">Completed</h3>
          </div>
          <p className="text-2xl font-bold text-green-600 mt-2">{completedTasks}</p>
        </div>
        <div className="rounded-lg border bg-card p-6">
          <div className="flex items-center gap-2">
            <Calendar className="h-5 w-5 text-red-600" />
            <h3 className="font-semibold">Overdue</h3>
          </div>
          <p className="text-2xl font-bold text-red-600 mt-2">{overdueTasks}</p>
        </div>
      </div>

      {/* Search */}
      <div className="relative max-w-sm">
        <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 h-4 w-4 text-muted-foreground" />
        <Input
          placeholder="Search tasks..."
          value={searchTerm}
          onChange={(e) => setSearchTerm(e.target.value)}
          className="pl-10"
        />
      </div>

      {/* Table */}
      <div className="rounded-md border border-border bg-card">
        <Table>
          <TableHeader>
            <TableRow className="hover:bg-muted/50">
              <TableHead className="text-card-foreground">Task Name</TableHead>
              <TableHead className="text-card-foreground">Description</TableHead>
              <TableHead className="text-card-foreground">Assigned To</TableHead>
              <TableHead className="text-card-foreground">Due Date</TableHead>
              <TableHead className="text-card-foreground">Status</TableHead>
              <TableHead className="w-[50px]"></TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {filteredTasks.length === 0 ? (
              <TableRow>
                <TableCell colSpan={6} className="text-center py-8 text-muted-foreground">
                  {searchTerm ? "No tasks found matching your search." : "No tasks yet. Create your first task!"}
                </TableCell>
              </TableRow>
            ) : (
              filteredTasks.map((task) => (
                <TableRow key={task.id} className="hover:bg-muted/50">
                  <TableCell className="font-medium text-card-foreground">
                    <div className="flex items-center gap-2">
                      <Button
                        variant="ghost"
                        size="sm"
                        onClick={() => handleToggleTaskStatus(task)}
                        className="p-1 h-auto"
                      >
                        {task.status === "Completed" ? (
                          <CheckSquare className="h-4 w-4 text-green-600" />
                        ) : (
                          <div className="h-4 w-4 border-2 border-gray-300 rounded" />
                        )}
                      </Button>
                      <span className={task.status === "Completed" ? "line-through text-muted-foreground" : ""}>
                        {task.task_name}
                      </span>
                    </div>
                  </TableCell>
                  <TableCell className="text-card-foreground max-w-xs truncate">
                    {task.task_description || "No description"}
                  </TableCell>
                  <TableCell className="text-card-foreground">
                    <div className="flex items-center gap-2">
                      <Calendar className="h-4 w-4 text-muted-foreground" />
                      {getUserName(task.assigned_to)}
                    </div>
                  </TableCell>
                  <TableCell className="text-card-foreground">
                    <div
                      className={`flex items-center gap-2 ${isOverdue(task.due_date, task.status) ? "text-red-600" : ""}`}
                    >
                      <Calendar className="h-4 w-4" />
                      {formatDate(task.due_date)}
                      {isOverdue(task.due_date, task.status) && (
                        <Badge variant="outline" className="bg-red-50 text-red-700 border-red-200 text-xs">
                          Overdue
                        </Badge>
                      )}
                    </div>
                  </TableCell>
                  <TableCell>{getStatusBadge(task.status)}</TableCell>
                  <TableCell>
                    <DropdownMenu>
                      <DropdownMenuTrigger asChild>
                        <Button variant="ghost" size="sm">
                          <MoreHorizontal className="h-4 w-4" />
                        </Button>
                      </DropdownMenuTrigger>
                      <DropdownMenuContent align="end">
                        <DropdownMenuItem onClick={() => openEditDialog(task)}>
                          <Edit className="mr-2 h-4 w-4" />
                          Edit
                        </DropdownMenuItem>
                        <DropdownMenuItem onClick={() => handleToggleTaskStatus(task)}>
                          <CheckSquare className="mr-2 h-4 w-4" />
                          {task.status === "Pending" ? "Mark Complete" : "Mark Pending"}
                        </DropdownMenuItem>
                        <DropdownMenuItem
                          className="text-destructive"
                          onClick={() => task.id && handleDeleteTask(task.id)}
                        >
                          <Trash2 className="mr-2 h-4 w-4" />
                          Delete
                        </DropdownMenuItem>
                      </DropdownMenuContent>
                    </DropdownMenu>
                  </TableCell>
                </TableRow>
              ))
            )}
          </TableBody>
        </Table>
      </div>

      {/* Edit Dialog */}
      <Dialog open={isEditDialogOpen} onOpenChange={setIsEditDialogOpen}>
        <DialogContent className="sm:max-w-[500px]">
          <DialogHeader>
            <DialogTitle>Edit Task</DialogTitle>
            <DialogDescription>Update the task information.</DialogDescription>
          </DialogHeader>
          <div className="grid gap-4 py-4">
            <div className="grid grid-cols-4 items-center gap-4">
              <Label htmlFor="edit_task_name" className="text-right">
                Task Name
              </Label>
              <Input
                id="edit_task_name"
                value={formData.task_name}
                onChange={(e) => setFormData((prev) => ({ ...prev, task_name: e.target.value }))}
                className="col-span-3"
                required
              />
            </div>
            <div className="grid grid-cols-4 items-center gap-4">
              <Label htmlFor="edit_task_description" className="text-right">
                Description
              </Label>
              <Textarea
                id="edit_task_description"
                value={formData.task_description}
                onChange={(e) => setFormData((prev) => ({ ...prev, task_description: e.target.value }))}
                className="col-span-3"
                rows={3}
              />
            </div>
            <div className="grid grid-cols-4 items-center gap-4">
              <Label htmlFor="edit_due_date" className="text-right">
                Due Date
              </Label>
              <Input
                id="edit_due_date"
                type="datetime-local"
                value={formData.due_date}
                onChange={(e) => setFormData((prev) => ({ ...prev, due_date: e.target.value }))}
                className="col-span-3"
                required
              />
            </div>
            <div className="grid grid-cols-4 items-center gap-4">
              <Label htmlFor="edit_assigned_to" className="text-right">
                Assigned To
              </Label>
              <Select
                value={formData.assigned_to}
                onValueChange={(value) => setFormData((prev) => ({ ...prev, assigned_to: value }))}
              >
                <SelectTrigger className="col-span-3">
                  <SelectValue placeholder="Select user" />
                </SelectTrigger>
                <SelectContent>
                  {users.map((user) => (
                    <SelectItem key={user.id} value={user.id!.toString()}>
                      {user.username}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>
            <div className="grid grid-cols-4 items-center gap-4">
              <Label htmlFor="edit_status" className="text-right">
                Status
              </Label>
              <Select
                value={formData.status}
                onValueChange={(value: "Pending" | "Completed") => setFormData((prev) => ({ ...prev, status: value }))}
              >
                <SelectTrigger className="col-span-3">
                  <SelectValue placeholder="Select status" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="Pending">Pending</SelectItem>
                  <SelectItem value="Completed">Completed</SelectItem>
                </SelectContent>
              </Select>
            </div>
          </div>
          <DialogFooter>
            <Button type="submit" onClick={handleEditTask} className="bg-cyan-600 hover:bg-cyan-700">
              Update Task
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  )
}
