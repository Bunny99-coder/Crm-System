"use client"

import { useState, useEffect } from "react"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Textarea } from "@/components/ui/textarea"
import { Badge } from "@/components/ui/badge"
import { Plus, Edit, Trash2, Save, X } from "lucide-react"
import { api } from "@/lib/api"

interface Note {
  id: number
  content: string
  created_at?: string
  updated_at?: string
  created_by: number
}

interface DealNotesProps {
  dealId: number
}

// Helper to safely format optional date strings
const formatDate = (dateStr?: string) => {
  if (!dateStr) return "-"
  const date = new Date(dateStr)
  if (isNaN(date.getTime())) return "-"
  return date.toLocaleDateString()
}

export function DealNotes({ dealId }: DealNotesProps) {
  const [notes, setNotes] = useState<Note[]>([])
  const [isLoading, setIsLoading] = useState(true)
  const [isCreating, setIsCreating] = useState(false)
  const [editingId, setEditingId] = useState<number | null>(null)
  const [newNote, setNewNote] = useState("")
  const [editContent, setEditContent] = useState("")

  useEffect(() => {
    loadNotes()
  }, [dealId])

  const loadNotes = async () => {
    try {
      setIsLoading(true)
      const notesData = await api.getNotesForDeal(dealId)
      setNotes(notesData || [])
    } catch (err) {
      console.error("Failed to load notes:", err)
      setNotes([])
    } finally {
      setIsLoading(false)
    }
  }

  const handleCreateNote = async () => {
    if (!newNote.trim()) return
    try {
      await api.createNoteForDeal(dealId, { content: newNote })
      setNewNote("")
      setIsCreating(false)
      loadNotes()
    } catch (err) {
      console.error("Failed to create note:", err)
    }
  }

  const handleUpdateNote = async (noteId: number) => {
    if (!editContent.trim()) return
    try {
      await api.updateNoteForDeal(dealId, noteId, { content: editContent })
      setEditingId(null)
      setEditContent("")
      loadNotes()
    } catch (err) {
      console.error("Failed to update note:", err)
    }
  }

  const handleDeleteNote = async (noteId: number) => {
    if (!confirm("Are you sure you want to delete this note?")) return
    try {
      await api.deleteNoteForDeal(dealId, noteId)
      loadNotes()
    } catch (err) {
      console.error("Failed to delete note:", err)
    }
  }

  const startEditing = (note: Note) => {
    setEditingId(note.id)
    setEditContent(note.content)
  }

  const cancelEditing = () => {
    setEditingId(null)
    setEditContent("")
  }

  if (isLoading) {
    return <div className="text-center py-8">Loading notes...</div>
  }

  return (
    <div className="space-y-4">
      {/* Create Note */}
      <Card>
        <CardHeader>
          <div className="flex items-center justify-between">
            <div>
              <CardTitle>Deal Notes</CardTitle>
              <CardDescription>Add and manage notes for this deal</CardDescription>
            </div>
            {!isCreating && (
              <Button onClick={() => setIsCreating(true)} size="sm" className="bg-cyan-600 hover:bg-cyan-700">
                <Plus className="mr-2 h-4 w-4" />
                Add Note
              </Button>
            )}
          </div>
        </CardHeader>

        {isCreating && (
          <CardContent className="space-y-4">
            <Textarea
              placeholder="Enter your note..."
              value={newNote}
              onChange={(e) => setNewNote(e.target.value)}
              rows={3}
            />
            <div className="flex gap-2">
              <Button onClick={handleCreateNote} size="sm" className="bg-cyan-600 hover:bg-cyan-700">
                <Save className="mr-2 h-4 w-4" />
                Save Note
              </Button>
              <Button
                onClick={() => {
                  setIsCreating(false)
                  setNewNote("")
                }}
                variant="outline"
                size="sm"
              >
                <X className="mr-2 h-4 w-4" />
                Cancel
              </Button>
            </div>
          </CardContent>
        )}
      </Card>

      {/* Notes List */}
      <div className="space-y-4">
        {!notes || notes.length === 0 ? (
          <Card>
            <CardContent className="text-center py-8">
              <p className="text-muted-foreground">No notes yet. Add your first note to get started.</p>
            </CardContent>
          </Card>
        ) : (
          notes.map((note) => (
            <Card key={note.id}>
              <CardContent className="pt-6">
                {editingId === note.id ? (
                  <div className="space-y-4">
                    <Textarea value={editContent} onChange={(e) => setEditContent(e.target.value)} rows={3} />
                    <div className="flex gap-2">
                      <Button
                        onClick={() => handleUpdateNote(note.id)}
                        size="sm"
                        className="bg-cyan-600 hover:bg-cyan-700"
                      >
                        <Save className="mr-2 h-4 w-4" />
                        Save
                      </Button>
                      <Button onClick={cancelEditing} variant="outline" size="sm">
                        <X className="mr-2 h-4 w-4" />
                        Cancel
                      </Button>
                    </div>
                  </div>
                ) : (
                  <div className="space-y-3">
                    <p className="text-sm leading-relaxed">{note.content}</p>
                    <div className="flex items-center justify-between">
                      <div className="flex items-center gap-2">
                        <Badge variant="outline" className="text-xs">
                          Created {formatDate(note.created_at)}
                        </Badge>
                        {note.updated_at && note.updated_at !== note.created_at && (
                          <Badge variant="outline" className="text-xs">
                            Updated {formatDate(note.updated_at)}
                          </Badge>
                        )}
                      </div>
                      <div className="flex gap-2">
                        <Button onClick={() => startEditing(note)} variant="ghost" size="sm">
                          <Edit className="h-4 w-4" />
                        </Button>
                        <Button
                          onClick={() => handleDeleteNote(note.id)}
                          variant="ghost"
                          size="sm"
                          className="text-destructive hover:text-destructive"
                        >
                          <Trash2 className="h-4 w-4" />
                        </Button>
                      </div>
                    </div>
                  </div>
                )}
              </CardContent>
            </Card>
          ))
        )}
      </div>
    </div>
  )
}
