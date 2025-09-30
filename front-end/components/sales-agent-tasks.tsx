import { useState, useEffect } from "react"
import { format } from "date-fns"
import { api, type Task } from "@/lib/api"
import { useAuth } from "@/lib/auth"
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table"
import { Badge } from "@/components/ui/badge"

export function SalesAgentTasks() {
  const [tasks, setTasks] = useState<Task[]>([])
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const { user, loading } = useAuth() // Destructure loading

  console.log("SalesAgentTasks component rendered. User:", user);

  useEffect(() => {
    console.log("useEffect in SalesAgentTasks running. User:", user);
    if (loading) { // Wait for auth to load
      console.log("Auth is still loading in SalesAgentTasks.");
      return;
    }
    const fetchTasks = async () => {
      if (!user?.id) {
        setError("User not authenticated.")
        setIsLoading(false)
        return
      }
      console.log("User ID is valid, attempting to fetch tasks.");
      try {
        setIsLoading(true)
        const userTasks = await api.getTasks(user.id)
        setTasks(userTasks || [])
      } catch (err) {
        setError("Failed to load tasks.")
        console.error("Failed to fetch tasks:", err)
      } finally {
        setIsLoading(false)
      }
    }
    fetchTasks()
  }, [user, loading]) // Add loading to dependency array

  if (isLoading) {
    return <div className="p-4 text-center text-muted-foreground">Loading tasks...</div>
  }

  if (error) {
    return <div className="p-4 text-center text-destructive">Error: {error}</div>
  }

  return (
    <div className="rounded-md border bg-card">
      <div className="p-4 border-b">
        <h2 className="text-lg font-semibold">My Tasks</h2>
      </div>
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead>Task Name</TableHead>
            <TableHead>Description</TableHead>
            <TableHead>Due Date</TableHead>
            <TableHead>Status</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {tasks.length === 0 ? (
            <TableRow>
              <TableCell colSpan={4} className="text-center py-8 text-muted-foreground">
                No tasks assigned.
              </TableCell>
            </TableRow>
          ) : (
            tasks.map((task) => (
              <TableRow key={task.id}>
                <TableCell className="font-medium">{task.task_name}</TableCell>
                <TableCell>{task.task_description}</TableCell>
                <TableCell>{format(new Date(task.due_date), "PPP")}</TableCell>
                <TableCell>
                  <Badge variant={task.status === "Completed" ? "default" : "secondary"}>
                    {task.status}
                  </Badge>
                </TableCell>
              </TableRow>
            ))
          )}
        </TableBody>
      </Table>
    </div>
  )
}
