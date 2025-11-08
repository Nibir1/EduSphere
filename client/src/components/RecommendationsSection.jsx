"use client";

import React, { useState, useEffect } from "react";
import {
  TrendingUp,
  Award,
  BookOpen,
  Download,
  Trash2,
  Loader2,
  Globe,
  FileText,
} from "lucide-react";
import api, { apiDownload } from "../api/axiosClient";

export default function RecommendationsSection({ uploadedDocuments }) {
  const [courses, setCourses] = useState([]);
  const [scholarships, setScholarships] = useState([]);
  const [summaries, setSummaries] = useState([]);
  const [aiSummary, setAiSummary] = useState("");
  const [lastRecoId, setLastRecoId] = useState(null);

  const [loading, setLoading] = useState(false);
  const [fetchingScholarships, setFetchingScholarships] = useState(false);
  const [generatingSummary, setGeneratingSummary] = useState(false);
  const [saving, setSaving] = useState(false);
  const [loadingSummaries, setLoadingSummaries] = useState(false);
  const [error, setError] = useState("");

  // Load last recommendation ID from localStorage
  useEffect(() => {
    const rid = localStorage.getItem("last_reco_id");
    if (rid) setLastRecoId(parseInt(rid, 10));
  }, []);

  // Fetch saved summaries
  const fetchSummaries = async () => {
    setLoadingSummaries(true);
    try {
      const res = await api.get("/summaries");
      setSummaries(res.data);
    } catch (err) {
      console.error("Failed to load summaries:", err);
    } finally {
      setLoadingSummaries(false);
    }
  };

  useEffect(() => {
    fetchSummaries();
  }, []);

  // 🔹 Generate AI-based course recommendations
  const fetchRecommendations = async () => {
    setLoading(true);
    setError("");
    try {
      const res = await api.post("/recommendations/generate");
      setCourses(res.data.courses || []);
      setScholarships(res.data.scholarships || []);
      if (res.data.id) localStorage.setItem("last_reco_id", res.data.id);
    } catch (err) {
      console.error("Failed to fetch recommendations:", err);
      setError(
        err.response?.data?.error ||
        "Failed to generate recommendations. Try again."
      );
    } finally {
      setLoading(false);
    }
  };

  // Automatically generate recommendations once
  useEffect(() => {
    fetchRecommendations();
  }, []);

  // 🔹 Generate AI transcript summary
  const generateSummary = async () => {
    if (generatingSummary || fetchingScholarships) return;
    setGeneratingSummary(true);
    try {
      const res = await api.post("/summaries/generate");
      setAiSummary(res.data.summary_text || res.data.text || "");
      alert("Summary generated successfully.");
    } catch (err) {
      console.error(err);
      alert(err.response?.data?.error || "Failed to generate summary.");
    } finally {
      setGeneratingSummary(false);
    }
  };

  // 🔹 Fetch scholarships
  const fetchScholarships = async () => {
    if (fetchingScholarships || generatingSummary) return;
    setFetchingScholarships(true);
    try {
      const res = await api.post("/scholarships/generate");
      const list = Array.isArray(res.data?.scholarships)
        ? res.data.scholarships
        : [];
      setScholarships(list);
      if (list.length === 0)
        alert("No scholarships found for this profile yet.");
    } catch (e) {
      console.error(e);
      alert(e.response?.data?.error || "Failed to fetch scholarships.");
    } finally {
      setFetchingScholarships(false);
    }
  };

  // 🔹 Save unified summary PDF (summary + scholarships + recommendations)
  const saveSummaryPDF = async () => {
    if (!lastRecoId) return alert("No recommendation available to save yet.");
    if (!aiSummary)
      return alert("Generate a transcript summary before saving.");
    setSaving(true);
    try {
      await api.post("/summaries", {
        recommendation_id: lastRecoId,
        summary_text: aiSummary,
        include_scholarships: scholarships.length > 0, // only include if user fetched
      });
      alert("Summary PDF saved (includes courses and scholarships).");
      await fetchSummaries();
    } catch (e) {
      console.error(e);
      alert(e.response?.data?.error || "Failed to save summary.");
    } finally {
      setSaving(false);
    }
  };

  // 🔹 Download PDF
  const handleDownload = async (id) => {
    try {
      await apiDownload(`/summaries/${id}/download`, `summary_${id}.pdf`);
    } catch (error) {
      console.error("PDF download failed:", error);
      alert("Failed to download PDF. Please try again.");
    }
  };

  // 🔹 Delete summary
  const handleDelete = async (id) => {
    if (!window.confirm("Are you sure you want to delete this summary?")) return;
    try {
      await api.delete(`/summaries/${id}`);
      alert("Summary deleted successfully.");
      await fetchSummaries();
    } catch (err) {
      console.error(err);
      alert("Failed to delete summary.");
    }
  };

  // 🔹 Loading and error states
  if (loading) {
    return (
      <div className="flex items-center justify-center py-20">
        <Loader2 className="h-6 w-6 animate-spin text-blue-600 mr-2" />
        <span className="text-gray-600">
          Generating personalized recommendations...
        </span>
      </div>
    );
  }

  if (error) {
    return (
      <div className="text-center text-red-600 py-10">
        ⚠️ {error}
        <br />
        <button
          onClick={fetchRecommendations}
          className="mt-3 text-sm text-blue-600 underline"
        >
          Retry
        </button>
      </div>
    );
  }

  // 🔹 Render UI
  return (
    <div className="space-y-8">
      {/* Summary Stats */}
      <div className="grid gap-4 md:grid-cols-3">
        <StatCard
          title="Documents Analyzed"
          value={uploadedDocuments.length}
          icon={<BookOpen className="h-6 w-6 text-blue-600" />}
        />
        <StatCard
          title="Courses Found"
          value={courses.length}
          icon={<TrendingUp className="h-6 w-6 text-indigo-600" />}
        />
        <StatCard
          title="Scholarships"
          value={scholarships.length}
          icon={<Award className="h-6 w-6 text-green-600" />}
        />
      </div>

      {/* Recommended Courses */}
      <SectionTitle>Recommended Courses</SectionTitle>
      <CourseList courses={courses} />

      {/* Scholarships */}
      <div className="flex items-center justify-between">
        <SectionTitle>Scholarship Opportunities</SectionTitle>
        <button
          onClick={fetchScholarships}
          disabled={fetchingScholarships || generatingSummary || saving}
          className="inline-flex items-center gap-2 rounded-lg bg-green-600 px-4 py-2 text-white font-semibold hover:bg-green-700 disabled:opacity-60"
        >
          <Globe className="w-4 h-4" />
          {fetchingScholarships ? "Searching..." : "Find Scholarships"}
        </button>
      </div>

      <ScholarshipList
        scholarships={scholarships}
        loading={fetchingScholarships}
      />

      {/* --- AI Summary Section --- */}
      <div className="p-6 rounded-lg border border-gray-300 bg-white">
        <h2 className="text-lg font-bold mb-3 flex items-center gap-2">
          <FileText className="w-5 h-5 text-indigo-600" /> Transcript Summary
        </h2>
        {aiSummary ? (
          <p className="text-gray-700 text-sm whitespace-pre-line mb-4">
            {aiSummary}
          </p>
        ) : (
          <p className="text-gray-500 text-sm">
            Generate a concise summary of your transcript using AI.
          </p>
        )}
        <div className="flex gap-3">
          <button
            onClick={generateSummary}
            disabled={
              generatingSummary || fetchingScholarships || saving || loading
            }
            className="rounded-lg bg-indigo-600 px-4 py-2 text-white font-semibold hover:bg-indigo-700 disabled:opacity-60"
          >
            {generatingSummary ? "Generating..." : "Generate Summary"}
          </button>

          {aiSummary && (
            <button
              onClick={saveSummaryPDF}
              disabled={
                saving ||
                fetchingScholarships || // 🔹 disable while fetching scholarships
                !aiSummary.trim()       // 🔹 disable if summary is empty
              }
              className="rounded-lg bg-blue-600 px-4 py-2 text-white font-semibold hover:bg-blue-700 disabled:opacity-60"
            >
              {saving
                ? "Saving..."
                : fetchingScholarships
                  ? "Please wait (loading scholarships)..."
                  : "Save Full Report (PDF)"}
            </button>
          )}

        </div>
      </div>

      {/* Saved Summaries */}
      <div>
        <SectionTitle>Saved Results</SectionTitle>
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
  );
}

// --- Small UI Components ---
const StatCard = ({ title, value, icon }) => (
  <div className="rounded-lg border border-gray-300 bg-white p-6">
    <div className="flex items-center justify-between">
      <div>
        <p className="text-sm text-gray-500">{title}</p>
        <p className="mt-1 text-3xl font-bold text-gray-900">{value}</p>
      </div>
      <div className="rounded-full bg-gray-100 p-3">{icon}</div>
    </div>
  </div>
);

const SectionTitle = ({ children }) => (
  <h2 className="mb-4 text-xl font-bold text-gray-900">{children}</h2>
);

const CourseList = ({ courses }) => (
  <div className="grid gap-4">
    {courses.length === 0 ? (
      <p className="text-gray-500 text-sm">
        No recommendations yet. Try uploading a transcript.
      </p>
    ) : (
      courses.map((course, idx) => (
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
                <span className="text-sm font-semibold text-blue-600">
                  {Math.round(course.match)}%
                </span>
              </div>
            </div>
          </div>
        </div>
      ))
    )}
  </div>
);

const ScholarshipList = ({ scholarships, loading }) => {
  if (loading)
    return (
      <div className="flex items-center justify-center py-8 text-gray-600">
        <Loader2 className="h-5 w-5 animate-spin mr-2 text-green-600" />
        Searching scholarships...
      </div>
    );
  return (
    <div className="grid gap-4">
      {scholarships.map((sch, idx) => (
        <div
          key={idx}
          className="rounded-lg border border-gray-300 bg-white p-6 hover:shadow-md transition-shadow"
        >
          <div className="flex items-start justify-between gap-4">
            <div className="flex-1">
              <h3 className="font-semibold text-gray-900">{sch.title}</h3>
              <p className="mt-1 text-sm text-gray-500">{sch.description}</p>
              {sch.link && (
                <a
                  href={sch.link}
                  target="_blank"
                  rel="noreferrer"
                  className="text-sm text-green-700 hover:underline inline-block mt-1"
                >
                  View details →
                </a>
              )}
            </div>
            <div className="flex flex-col items-end gap-2">
              <div className="rounded-full bg-green-100 px-3 py-1">
                <span className="text-sm font-semibold text-green-600">
                  {Math.round(sch.match)}%
                </span>
              </div>
            </div>
          </div>
        </div>
      ))}
    </div>
  );
};
