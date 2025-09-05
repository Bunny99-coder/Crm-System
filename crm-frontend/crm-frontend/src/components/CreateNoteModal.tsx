// File: src/components/CreateNoteModal.tsx
import { Modal, Button, Stack, Textarea } from '@mantine/core';
import { useForm, isNotEmpty } from '@mantine/form';
import { useEffect } from 'react';
import { createNote, updateNote } from '../services/apiService';
import type { Note, CreateNotePayload, UpdateNotePayload } from '../services/apiService';
import { useAuth } from '../context/AuthContext';
import { getUserIdFromToken } from '../util/auth'; // We'll create this helper

interface CreateNoteModalProps {
  opened: boolean;
  onClose: () => void;
  onSuccess: () => void;
  noteToEdit: Note | null;
}

const CreateNoteModal = ({ opened, onClose, onSuccess, noteToEdit }: CreateNoteModalProps) => {
  const isEditing = !!noteToEdit;
  const { token } = useAuth();
  
  const form = useForm({
    initialValues: { note_text: '' },
    validate: { note_text: isNotEmpty('Note cannot be empty') },
  });

  useEffect(() => {
    if (isEditing) {
      form.setValues({ note_text: noteToEdit.note_text });
    } else {
      form.reset();
    }
  }, [noteToEdit, opened]);

  const handleSubmit = async (values: typeof form.values) => {
    try {
      if (isEditing) {
        const payload: UpdateNotePayload = {
          note_text: values.note_text,
          note_date: new Date().toISOString(), // Update the date on edit
        };
        await updateNote(noteToEdit.id, payload);
      } else {
        const payload: CreateNotePayload = {
          note_text: values.note_text,
          user_id: getUserIdFromToken(token) || 0,
        };
        await createNote(payload);
      }
      onSuccess();
      onClose();
    } catch (error) {
      console.error("Failed to save note", error);
    }
  };

  return (
    <Modal opened={opened} onClose={onClose} title={isEditing ? 'Edit Note' : 'Create New Note'} centered>
      <form onSubmit={form.onSubmit(handleSubmit)}>
        <Stack>
          <Textarea label="Note" required minRows={5} {...form.getInputProps('note_text')} />
          <Button type="submit" mt="md">{isEditing ? 'Update Note' : 'Create Note'}</Button>
        </Stack>
      </form>
    </Modal>
  );
};

export default CreateNoteModal;