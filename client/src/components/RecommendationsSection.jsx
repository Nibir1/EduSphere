"use client"

import React, { useState, useEffect } from "react"
import { TrendingUp, Award, BookOpen, Download, Trash2 } from "lucide-react"
import api from "../api/axiosClient"

export default function RecommendationsSection({ recommendations, uploadedDocuments }) {
  const [saving, setSaving] = useState(false)
  const [loadingSummaries, setLoadingSummaries] = useState(false)
  const [summaries, setSummaries] = useState([])
  const [lastRecoId, setLastRecoId] = useState(null)

  // Retrieve last recommendation ID from localStorage (set in UploadSection)
  useEffect(() => {
    const rid = localStorage.getItem("last_reco_id")
    if (rid) setLastRecoId(parseInt(rid, 10))
  }, [])

  // Fetch list of saved summaries
  const fetchSummaries = async () => {
    setLoadingSummaries(true)
    try {
      const res = await api.get("/summaries")
      setSummaries(res.data)
    } catch (err) {
      console.error(err)
    } finally {
      setLoadingSummaries(false)
    }
  }

  useEffect(() => {
    fetchSummaries()
  }, [])

  // Save the current recommendation as a PDF summary
  const saveSummary = async () => {
    if (!lastRecoId) return alert("No recommendation available to save yet.")
    setSaving(true)
    try {
      const res = await api.post("/summaries", { recommendation_id: lastRecoId })
      const pdfPath = res.data.pdf_path
      alert("✅ Summary PDF saved successfully!")
      await fetchSummaries()
    } catch (e) {
      console.error(e)
      alert(e.response?.data?.error || "Failed to save summary.")
    } finally {
      setSaving(false)
    }
  }

  // Download a summary
  const handleDownload = async (id) => {
    try {
      const response = await api.get(`/summaries/${id}/download`, { responseType: "blob" })
      const blob = new Blob([response.data], { type: "application/pdf" })
      const url = window.URL.createObjectURL(blob)
      const link = document.createElement("a")
      link.href = url
      link.setAttribute("download", `summary_${id}.pdf`)
      document.body.appendChild(link)
      link.click()
      link.parentNode.removeChild(link)
      window.URL.revokeObjectURL(url)
    } catch (error) {
      console.error(error)
      alert("Failed to download PDF")
    }
  }

  // Delete a summary
  const handleDelete = async (id) => {
    if (!window.confirm("Are you sure you want to delete this summary?")) return
    try {
      await api.delete(`/summaries/${id}`)
      alert("🗑️ Summary deleted successfully")
      await fetchSummaries()
    } catch (err) {
      console.error(err)
      alert("Failed to delete summary.")
    }
  }

  const courses = recommendations.filter((r) => r.type === "course")
  const scholarships = recommendations.filter((r) => r.type === "scholarship")

  return (
    <div className="space-y-8">
      {/* Summary Stats */}
      <div className="grid gap-4 md:grid-cols-3">
        <div className="rounded-lg border border-gray-300 bg-white p-6">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm text-gray-500">Documents Analyzed</p>
              <p className="mt-1 text-3xl font-bold text-gray-900">{uploadedDocuments.length}</p>
            </div>
            <div className="rounded-full bg-blue-100 p-3">
              <BookOpen className="h-6 w-6 text-blue-600" />
            </div>
          </div>
        </div>

        <div className="rounded-lg border border-gray-300 bg-white p-6">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm text-gray-500">Courses Found</p>
              <p className="mt-1 text-3xl font-bold text-gray-900">{courses.length}</p>
            </div>
            <div className="rounded-full bg-indigo-100 p-3">
              <TrendingUp className="h-6 w-6 text-indigo-600" />
            </div>
          </div>
        </div>

        <div className="rounded-lg border border-gray-300 bg-white p-6">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm text-gray-500">Scholarships</p>
              <p className="mt-1 text-3xl font-bold text-gray-900">{scholarships.length}</p>
            </div>
            <div className="rounded-full bg-green-100 p-3">
              <Award className="h-6 w-6 text-green-600" />
            </div>
          </div>
        </div>
      </div>

      {/* Recommended Courses */}
      <div>
        <h2 className="mb-4 text-xl font-bold text-gray-900">Recommended Courses</h2>
        <div className="grid gap-4">
          {courses.length === 0 && (
            <p className="text-gray-500 text-sm">No recommendations yet. Try uploading a transcript.</p>
          )}
          {courses.map((course, idx) => (
            <div
              key={idx}
              className="rounded-lg border border-gray-300 bg-white p-6 hover:shadow-md transition-shadow"
            >
              <div className="flex items-start justify-between gap-4">
                <div className="flex-1">
                  <h3 className="font-semibold text-gray-900">{course.title}</h3>
                  <p className="mt-1 text-sm text-gray-500">{course.description}</p>
                </div>
                <div className="flex flex-col items-end gap-2">
                  <div className="rounded-full bg-blue-100 px-3 py-1">
                    <span className="text-sm font-semibold text-blue-600">{course.match}%</span>
                  </div>
                  <button className="text-xs font-medium text-blue-600 hover:underline">Learn More →</button>
                </div>
              </div>
            </div>
          ))}
        </div>
      </div>

      {/* Recommended Scholarships */}
      <div>
        <h2 className="mb-4 text-xl font-bold text-gray-900">Scholarship Opportunities</h2>
        <div className="grid gap-4">
          {scholarships.length === 0 && (
            <p className="text-gray-500 text-sm">No scholarships found for this profile yet.</p>
          )}
          {scholarships.map((scholarship, idx) => (
            <div
              key={idx}
              className="rounded-lg border border-gray-300 bg-white p-6 hover:shadow-md transition-shadow"
            >
              <div className="flex items-start justify-between gap-4">
                <div className="flex-1">
                  <h3 className="font-semibold text-gray-900">{scholarship.title}</h3>
                  <p className="mt-1 text-sm text-gray-500">{scholarship.description}</p>
                </div>
                <div className="flex flex-col items-end gap-2">
                  <div className="rounded-full bg-green-100 px-3 py-1">
                    <span className="text-sm font-semibold text-green-600">{scholarship.match}%</span>
                  </div>
                  <button className="text-xs font-medium text-green-600 hover:underline">Apply Now →</button>
                </div>
              </div>
            </div>
          ))}
        </div>
      </div>

      {/* Save Result Button */}
      {lastRecoId && (
        <div className="flex justify-end">
          <button
            onClick={saveSummary}
            disabled={saving}
            className="rounded-lg bg-indigo-600 px-5 py-2.5 text-white font-semibold hover:bg-indigo-700 disabled:opacity-60 transition-all"
          >
            {saving ? "Saving..." : "Save Result (PDF)"}
          </button>
        </div>
      )}

      {/* Saved Summaries */}
      <div>
        <h2 className="mb-4 text-xl font-bold text-gray-900">Saved Results</h2>
        {loadingSummaries ? (
          <p className="text-gray-500 text-sm">Loading saved summaries...</p>
        ) : summaries.length === 0 ? (
          <p className="text-gray-500 text-sm">No saved summaries yet.</p>
        ) : (
          <div className="grid gap-3">
            {summaries.map((s) => (
              <div
                key={s.id}
                className="flex items-center justify-between p-4 rounded-lg border border-gray-300 bg-white hover:bg-gray-50 transition"
              >
                <div>
                  <p className="font-medium text-gray-900">Summary #{s.id}</p>
                  <p className="text-xs text-gray-500">
                    Created: {new Date(s.created_at).toLocaleString()}
                  </p>
                </div>
                <div className="flex gap-3">
                  <button
                    onClick={() => handleDownload(s.id)}
                    className="flex items-center gap-1 text-blue-600 hover:text-blue-800 transition"
                  >
                    <Download className="w-4 h-4" /> Download
                  </button>
                  <button
                    onClick={() => handleDelete(s.id)}
                    className="flex items-center gap-1 text-red-600 hover:text-red-800 transition"
                  >
                    <Trash2 className="w-4 h-4" /> Delete
                  </button>
                </div>
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  )
}
