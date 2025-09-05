// File: src/components/CreateCommLogModal.tsx
import { Modal, Button, Select, Stack, Textarea } from '@mantine/core';
import { useForm, isNotEmpty } from '@mantine/form';
import { DateTimePicker } from '@mantine/dates';
import { useEffect } from 'react';
import { createCommLog, updateCommLog } from '../services/apiService';
import type { CommLog, CreateCommLogPayload } from '../services/apiService';
import { useAuth } from '../context/AuthContext';

interface CreateCommLogModalProps {
  opened: boolean;
  onClose: () => void;
  onSuccess: () => void;
  logToEdit: CommLog | null;
  contactId: number; // We need to know which contact this log is for
}

const CreateCommLogModal = ({ opened, onClose, onSuccess, logToEdit, contactId }: CreateCommLogModalProps) => {
  const isEditing = !!logToEdit;
  const { token } = useAuth(); // We'll need the user's ID from the token

  // A simple helper to decode the user ID from the JWT
  const getUserIdFromToken = () => {
    if (!token) return 0;
    try {
      const payload = JSON.parse(atob(token.split('.')[1]));
      return payload.user_id || 0;
    } catch (e) {
      return 0;
    }
  };
  
  const interactionTypes = ['Call', 'Email', 'Meeting', 'Site Visit'];

  const form = useForm({
    initialValues: {
      interaction_date: new Date(),
      interaction_type: 'Call',
      notes: '',
    },
    validate: {
      interaction_type: isNotEmpty('Interaction type is required'),
    },
  });

  useEffect(() => {
    if (isEditing) {
      form.setValues({
        interaction_date: new Date(logToEdit.interaction_date),
        interaction_type: logToEdit.interaction_type,
        notes: logToEdit.notes || '',
      });
    } else {
      form.reset();
    }
  }, [logToEdit, opened]);

  const handleSubmit = async (values: typeof form.values) => {
    try {
      const payload = {
        interaction_date: values.interaction_date.toISOString(),
        interaction_type: values.interaction_type,
        notes: values.notes,
        // The rest of the payload depends on create vs update
      };

      if (isEditing) {
        await updateCommLog(logToEdit.id, payload);
      } else {
        const createPayload: CreateCommLogPayload = {
          ...payload,
          contact_id: contactId,
          user_id: getUserIdFromToken(), // Get the logged-in user's ID
        };
        await createCommLog(createPayload);
      }
      onSuccess();
      onClose();
    } catch (error) {
      console.error("Failed to save communication log", error);
    }
  };

  return (
    <Modal opened={opened} onClose={onClose} title={isEditing ? 'Edit Log Entry' : 'Add New Log Entry'} centered>
      <form onSubmit={form.onSubmit(handleSubmit)}>
        <Stack>
          <DateTimePicker
            label="Interaction Date & Time"
            placeholder="Pick date and time"
            required
            {...form.getInputProps('interaction_date')}
          />
          <Select
            label="Interaction Type"
            data={interactionTypes}
            required
            {...form.getInputProps('interaction_type')}
          />
          <Textarea
            label="Notes"
            placeholder="Details of the interaction..."
            minRows={4}
            {...form.getInputProps('notes')}
          />
          <Button type="submit" mt="md">{isEditing ? 'Update Log' : 'Add Log'}</Button>
        </Stack>
      </form>
    </Modal>
  );
};

export default CreateCommLogModal;