// Replace the contents of CreateContactModal.tsx

import { Modal, Button, TextInput, Stack, Title } from '@mantine/core';
import { useForm } from '@mantine/form';
import { useEffect } from 'react';
import { createContact, updateContact } from '../services/apiService';
import type { Contact, CreateContactPayload } from '../services/apiService';

interface CreateContactModalProps {
  opened: boolean;
  onClose: () => void;
  onSuccess: () => void;
  contactToEdit: Contact | null; // <-- NEW: Pass the contact to edit
}

const CreateContactModal = ({ opened, onClose, onSuccess, contactToEdit }: CreateContactModalProps) => {
  const isEditing = !!contactToEdit; // Check if we are in "edit" mode

  const form = useForm<CreateContactPayload>({
    initialValues: {
      first_name: '',
      last_name: '',
      email: '',
      primary_phone: '',
    },
    validate: {
      first_name: (value) => (value.trim().length > 0 ? null : 'First name is required'),
      primary_phone: (value) => (value.trim().length > 0 ? null : 'Primary phone is required'),
    },
  });

  // This effect will run when the modal opens or when contactToEdit changes.
  // It pre-fills the form if we are in edit mode.
  useEffect(() => {
    if (isEditing) {
      form.setValues({
        first_name: contactToEdit.first_name,
        last_name: contactToEdit.last_name,
        email: contactToEdit.email || '',
        primary_phone: contactToEdit.primary_phone,
      });
    } else {
      form.reset(); // Clear the form if we are in create mode
    }
  }, [contactToEdit, opened]);

  const handleSubmit = async (values: CreateContactPayload) => {
    try {
      if (isEditing) {
        await updateContact(contactToEdit.id, values);
      } else {
        await createContact(values);
      }
      onSuccess();
      onClose();
    } catch (error) {
      console.error("Failed to save contact", error);
    }
  };

  return (
    <Modal opened={opened} onClose={onClose} title={isEditing ? 'Edit Contact' : 'Create New Contact'}>
      <form onSubmit={form.onSubmit(handleSubmit)}>
        <Stack>
          <TextInput label="First Name" required {...form.getInputProps('first_name')} />
          <TextInput label="Last Name" {...form.getInputProps('last_name')} />
          <TextInput label="Email" type="email" {...form.getInputProps('email')} />
          <TextInput label="Primary Phone" required {...form.getInputProps('primary_phone')} />
          <Button type="submit" mt="md">{isEditing ? 'Update Contact' : 'Create Contact'}</Button>
        </Stack>
      </form>
    </Modal>
  );
};

export default CreateContactModal;