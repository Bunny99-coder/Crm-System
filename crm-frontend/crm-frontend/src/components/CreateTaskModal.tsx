// Replace the contents of src/components/CreateTaskModal.tsx

import { Modal, Button, TextInput, Select, Stack, Textarea } from '@mantine/core';
import { useForm, isNotEmpty } from '@mantine/form';
import { DatePickerInput } from '@mantine/dates';
import { useState, useEffect } from 'react';
import { createTask, updateTask, getUsers } from '../services/apiService';
import type { Task, CreateTaskPayload, UserSelectItem } from '../services/apiService';

interface CreateTaskModalProps {
  opened: boolean;
  onClose: () => void;
  onSuccess: () => void;
  taskToEdit: Task | null;
}

const CreateTaskModal = ({ opened, onClose, onSuccess, taskToEdit }: CreateTaskModalProps) => {
  const isEditing = !!taskToEdit;
  const [users, setUsers] = useState<{ value: string; label: string }[]>([]);
  // We need a loading state for the button to prevent double-clicks
  const [isSubmitting, setIsSubmitting] = useState(false);

  const form = useForm({
    initialValues: {
      task_name: '',
      task_description: '',
      due_date: new Date(),
      status: 'Pending',
      assigned_to: '',
    },
    validate: {
      task_name: isNotEmpty('Task name is required'),
      assigned_to: isNotEmpty('A user must be assigned'),
      due_date: (value) => (value ? null : 'Due date is required'),
    },
  });

  useEffect(() => {
    if (opened) {
      getUsers().then(usersData => {
        setUsers(usersData.map((u: UserSelectItem) => ({ value: String(u.id), label: u.username })));
      });

      if (isEditing) {
        form.setValues({
          task_name: taskToEdit.task_name,
          task_description: taskToEdit.task_description || '',
          due_date: new Date(taskToEdit.due_date),
          status: taskToEdit.status,
          assigned_to: String(taskToEdit.assigned_to),
        });
      } else {
        form.reset();
      }
    }
  }, [taskToEdit, opened]);

  // --- THIS IS THE CORRECTED SUBMIT FUNCTION ---
  const handleSubmit = async (values: typeof form.values) => {
    setIsSubmitting(true); // Disable the button
    try {
      // 1. Build the payload with the correct types
      const payload: CreateTaskPayload = {
        task_name: values.task_name,
        task_description: values.task_description,
        due_date: values.due_date.toISOString(), // Convert Date object to ISO string
        status: values.status,
        assigned_to: Number(values.assigned_to), // Convert string ID to number
      };

      // 2. Call the correct API function based on the mode
      if (isEditing) {
        await updateTask(taskToEdit.id, payload);
      } else {
        await createTask(payload);
      }
      
      // 3. Signal success to the parent page so it can refetch data
      onSuccess();
      // 4. Close the modal
      onClose();

    } catch (error) {
      console.error("Failed to save task", error);
      // In a real app, you would show an error notification here
      alert("Error: Could not save the task. Check the console for details.");
    } finally {
      setIsSubmitting(false); // Re-enable the button
    }
  };

  return (
    <Modal opened={opened} onClose={onClose} title={isEditing ? 'Edit Task' : 'Create New Task'} centered>
      <form onSubmit={form.onSubmit(handleSubmit)}>
        <Stack>
          <TextInput label="Task Name" required {...form.getInputProps('task_name')} />
          <Textarea label="Description" {...form.getInputProps('task_description')} />
          <DatePickerInput label="Due Date" required {...form.getInputProps('due_date')} />
          <Select label="Status" data={['Pending', 'Completed']} required {...form.getInputProps('status')} />
          <Select label="Assigned To" data={users} searchable required {...form.getInputProps('assigned_to')} />
          <Button type="submit" mt="md" loading={isSubmitting}>{isEditing ? 'Update Task' : 'Create Task'}</Button>
        </Stack>
      </form>
    </Modal>
  );
};

export default CreateTaskModal;