// Replace the entire contents of src/pages/TasksPage.tsx with this.

import { useState, useEffect, useCallback } from 'react';
import { getTasks, deleteTask } from '../services/apiService';
import type { Task } from '../services/apiService';
import { Table, Loader, Alert, Title, Button, Group, ActionIcon, Badge } from '@mantine/core';
import { IconPencil, IconTrash, IconAlertCircle } from '@tabler/icons-react';
import { useDisclosure } from '@mantine/hooks';
import CreateTaskModal from '../components/CreateTaskModal';
import { useData } from '../context/DataContext';

const TasksPage = () => {
  const [tasks, setTasks] = useState<Task[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [modalOpened, { open: openModal, close: closeModal }] = useDisclosure(false);
  const [taskToEdit, setTaskToEdit] = useState<Task | null>(null);

  const { userMap, loading: dataLoading } = useData();

  const fetchTasks = useCallback(async () => {
    try {
      setLoading(true);
      const data = await getTasks();
      setTasks(data || []);
      setError(null);
    } catch (err) {
      setError('Failed to load tasks. Are you logged in?');
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    fetchTasks();
  }, [fetchTasks]);

  // --- THIS IS THE MISSING LOGIC ---
  const handleEditClick = (task: Task) => {
    setTaskToEdit(task);
    openModal();
  };

  const handleCreateClick = () => {
    setTaskToEdit(null); // Make sure we are in "create" mode
    openModal(); // This function opens the modal
  };
  
  const handleDeleteClick = async (id: number) => {
    if (window.confirm('Are you sure you want to delete this task?')) {
      try {
        await deleteTask(id);
        fetchTasks();
      } catch (err) {
        setError('Failed to delete task.');
      }
    }
  };
  // ------------------------------------

  if ((loading || dataLoading) && tasks.length === 0) return <Loader color="blue" />;
  if (error) return <Alert color="red" title="Error" icon={<IconAlertCircle />}>{error}</Alert>;

  const rows = tasks.map((task) => (
    <Table.Tr key={task.id}>
      <Table.Td>{task.task_name}</Table.Td>
      <Table.Td>
        <Badge color={task.status === 'Completed' ? 'teal' : 'yellow'}>
          {task.status}
        </Badge>
      </Table.Td>
      <Table.Td>{new Date(task.due_date).toLocaleDateString()}</Table.Td>
      <Table.Td>{userMap.get(task.assigned_to) || `User ID: ${task.assigned_to}`}</Table.Td>
      <Table.Td>
        <Group gap="xs">
          <ActionIcon variant="subtle" onClick={() => handleEditClick(task)}><IconPencil size={16} /></ActionIcon>
          <ActionIcon variant="subtle" color="red" onClick={() => handleDeleteClick(task.id)}><IconTrash size={16} /></ActionIcon>
        </Group>
      </Table.Td>
    </Table.Tr>
  ));

  return (
    <>
      <CreateTaskModal
        opened={modalOpened}
        onClose={closeModal}
        onSuccess={fetchTasks}
        taskToEdit={taskToEdit}
      />

      <Group justify="space-between" mb="md">
        <Title order={2}>Tasks</Title>
        <Button onClick={handleCreateClick}>Create New Task</Button>
      </Group>

      <Table striped withTableBorder withColumnBorders>
        <Table.Thead>
          <Table.Tr>
            <Table.Th>Task</Table.Th>
            <Table.Th>Status</Table.Th>
            <Table.Th>Due Date</Table.Th>
            <Table.Th>Assigned To</Table.Th>
            <Table.Th>Actions</Table.Th>
          </Table.Tr>
        </Table.Thead>
        <Table.Tbody>
          {rows.length > 0 ? rows : (
            <Table.Tr><Table.Td colSpan={5} align="center">No tasks found. Create one!</Table.Td></Table.Tr>
          )}
        </Table.Tbody>
      </Table>
    </>
  );
};

export default TasksPage;