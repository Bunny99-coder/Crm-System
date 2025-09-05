// Replace the entire contents of your src/components/CreateEventModal.tsx file with this.

import { Modal, Button, TextInput, Stack, Textarea } from '@mantine/core';
import { useForm, isNotEmpty } from '@mantine/form';
import { DateTimePicker } from '@mantine/dates';
import { useEffect, useState } from 'react';
import { createEvent, updateEvent } from '../services/apiService';
import type { Event, CreateEventPayload, UpdateEventPayload } from '../services/apiService';
import { useAuth } from '../context/AuthContext';
import { getUserIdFromToken } from '../util/auth';

interface CreateEventModalProps {
  opened: boolean;
  onClose: () => void;
  onSuccess: () => void;
  eventToEdit: Event | null;
}

const CreateEventModal = ({ opened, onClose, onSuccess, eventToEdit }: CreateEventModalProps) => {
  const isEditing = !!eventToEdit;
  const { token } = useAuth();
  const [isSubmitting, setIsSubmitting] = useState(false);

  const form = useForm({
    initialValues: {
      event_name: '',
      event_description: '',
      start_time: new Date(),
      end_time: new Date(new Date().getTime() + 60 * 60 * 1000), // Default to 1 hour after now
      location: '',
    },
    validate: {
      event_name: isNotEmpty('Event name is required'),
      start_time: (value) => (value ? null : 'Start time is required'),
      end_time: (value, values) => 
        (value && values.start_time && value > values.start_time ? null : 'End time must be after start time'),
    },
  });

  useEffect(() => {
    if (opened) {
      if (isEditing) {
        form.setValues({
          event_name: eventToEdit.event_name,
          event_description: eventToEdit.event_description || '',
          start_time: new Date(eventToEdit.start_time),
          end_time: new Date(eventToEdit.end_time),
          location: eventToEdit.location || '',
        });
      } else {
        // Reset and set a sensible default for new events
        form.reset();
        const now = new Date();
        form.setFieldValue('start_time', now);
        form.setFieldValue('end_time', new Date(now.getTime() + 60 * 60 * 1000));
      }
    }
  }, [eventToEdit, opened]);


  // --- THIS IS THE CORRECTED SUBMIT FUNCTION ---
  const handleSubmit = async (values: typeof form.values) => {
    setIsSubmitting(true);
    try {
      // Create new Date objects from the form values to guarantee they are Date objects
      const startTime = new Date(values.start_time);
      const endTime = new Date(values.end_time);

      // Build the payload that will be sent to the API
      const basePayload = {
        event_name: values.event_name,
        event_description: values.event_description,
        start_time: startTime.toISOString(), // Now it's safe to call .toISOString()
        end_time: endTime.toISOString(),     // This is also safe
        location: values.location,
      };

      // Call the correct API function based on whether we are editing or creating
      if (isEditing) {
        const payload: UpdateEventPayload = {
          ...basePayload,
          organizer_id: eventToEdit.organizer_id,
        };
        await updateEvent(eventToEdit.id, payload);
      } else {
        const payload: CreateEventPayload = {
          ...basePayload,
          organizer_id: getUserIdFromToken(token) || 0,
        };
        await createEvent(payload);
      }

      onSuccess(); // Tell the parent page to refresh its data
      onClose();   // Close the modal

    } catch (error) {
      console.error("Failed to save event", error);
      alert("Error: Could not save the event. See console for details.");
    } finally {
      setIsSubmitting(false); // Re-enable the submit button
    }
  };

  return (
    <Modal opened={opened} onClose={onClose} title={isEditing ? 'Edit Event' : 'Schedule New Event'} centered>
      <form onSubmit={form.onSubmit(handleSubmit)}>
        <Stack>
          <TextInput label="Event Name" placeholder="e.g., Property Viewing" required {...form.getInputProps('event_name')} />
          <Textarea label="Description" placeholder="e.g., With Mr. & Mrs. Smith for Unit 101" {...form.getInputProps('event_description')} />
          <DateTimePicker label="Start Time" required {...form.getInputProps('start_time')} />
          <DateTimePicker label="End Time" required {...form.getInputProps('end_time')} />
          <TextInput label="Location" placeholder="e.g., 123 Main St, Downtown Condos" {...form.getInputProps('location')} />
          <Button type="submit" mt="md" loading={isSubmitting}>{isEditing ? 'Update Event' : 'Schedule Event'}</Button>
        </Stack>
      </form>
    </Modal>
  );
};

export default CreateEventModal;