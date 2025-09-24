"use client"

import { useState, useEffect } from "react"
import { Card, CardContent } from "@/components/ui/card"
import { Input } from "@/components/ui/input"
import { Button } from "@/components/ui/button"
import { api, type Note } from "@/lib/api"

interface ContactNotesProps {
  contactId: number
}

export default function ContactNotes({ contactId }: ContactNotesProps) {
  const [notes, setNotes] = useState<Note[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [newNoteContent, setNewNoteContent] = useState("")
  const [editingNoteId, setEditingNoteId] = useState<number | null>(null)
  const [editingContent, setEditingContent] = useState("")

  // Fetch notes
  useEffect(() => {
    async function fetchNotes() {
      try {
        const data = await api.getContactNotes(contactId)
        setNotes(data ?? [])
      } catch (err: any) {
        setError(err.message || "Something went wrong")
      } finally {
        setLoading(false)
      }
    }

    fetchNotes()
  }, [contactId])

  // Create new note
  async function handleCreateNote() {
    if (!newNoteContent.trim()) return
    try {
      const createdNote = await api.createContactNote(contactId, { content: newNoteContent })
      setNotes((prev) => [...prev, createdNote])
      setNewNoteContent("")
    } catch (err: any) {
      setError(err.message)
    }
  }

  // Delete note
  async function handleDeleteNote(id: number) {
    try {
      await api.deleteContactNote(contactId, id)
      setNotes((prev) => prev.filter((n) => n.id !== id))
    } catch (err: any) {
      setError(err.message)
    }
  }

  // Edit note
  function startEditNote(note: Note) {
    setEditingNoteId(note.id)
    setEditingContent(note.content)
  }

  function cancelEdit() {
    setEditingNoteId(null)
    setEditingContent("")
  }

  async function handleUpdateNote(id: number) {
    if (!editingContent.trim()) return
    try {
      const updatedNote = await api.updateContactNote(contactId, id, { content: editingContent })
      setNotes((prev) =>
        prev.map((n) => (n.id === id ? updatedNote : n))
      )
      cancelEdit()
    } catch (err: any) {
      setError(err.message)
    }
  }

  if (loading) return <p>Loading notes...</p>
  if (error) return <p className="text-red-500">{error}</p>

  return (
    <div className="space-y-4">
      {/* Error display */}
      {error && (
        <div className="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded">
          {error}
        </div>
      )}

      {/* New Note Input */}
      <div className="flex gap-2">
        <Input
          placeholder="Write a new note..."
          value={newNoteContent}
          onChange={(e) => setNewNoteContent(e.target.value)}
          onKeyPress={(e) => {
            if (e.key === 'Enter') handleCreateNote()
          }}
        />
        <Button onClick={handleCreateNote} disabled={!newNoteContent.trim()}>
          Add
        </Button>
      </div>

      {/* Notes List */}
      {notes.length === 0 ? (
        <Card>
          <CardContent className="text-center py-8">
            <p className="text-muted-foreground">No notes found.</p>
          </CardContent>
        </Card>
      ) : (
        <div className="space-y-3">
          {notes.map((note) => (
            <Card key={note.id} className="p-4">
              <CardContent className="p-0">
                {editingNoteId === note.id ? (
                  <div className="space-y-3">
                    <Input
                      value={editingContent}
                      onChange={(e) => setEditingContent(e.target.value)}
                      onKeyPress={(e) => {
                        if (e.key === 'Enter') handleUpdateNote(note.id)
                      }}
                    />
                    <div className="flex gap-2">
                      <Button 
                        size="sm" 
                        onClick={() => handleUpdateNote(note.id)}
                        disabled={!editingContent.trim()}
                      >
                        Save
                      </Button>
                      <Button size="sm" variant="outline" onClick={cancelEdit}>
                        Cancel
                      </Button>
                    </div>
                  </div>
                ) : (
                  <div className="space-y-2">
                    <div className="flex justify-between items-start">
                      <p className="text-sm">{note.content}</p>
                      <div className="flex gap-2">
                        <Button 
                          size="sm" 
                          variant="outline" 
                          onClick={() => startEditNote(note)}
                        >
                          Edit
                        </Button>
                        <Button 
                          size="sm" 
                          variant="destructive" 
                          onClick={() => handleDeleteNote(note.id)}
                        >
                          Delete
                        </Button>
                      </div>
                    </div>
                    <p className="text-xs text-muted-foreground">
                      {new Date(note.created_at).toLocaleString()}
                      {note.updated_at !== note.created_at && (
                        <span> • Edited: {new Date(note.updated_at).toLocaleString()}</span>
                      )}
                    </p>
                  </div>
                )}
              </CardContent>
            </Card>
          ))}
        </div>
      )}
    </div>
  )
}