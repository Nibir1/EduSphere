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
} from "lucide-react";
import api from "../api/axiosClient";

export default function RecommendationsSection({ uploadedDocuments }) {
  const [courses, setCourses] = useState([]);
  const [scholarships, setScholarships] = useState([]);
  const [loading, setLoading] = useState(false);
  const [fetchingScholarships, setFetchingScholarships] = useState(false);
  const [error, setError] = useState("");
  const [saving, setSaving] = useState(false);
  const [loadingSummaries, setLoadingSummaries] = useState(false);
  const [summaries, setSummaries] = useState([]);
  const [lastRecoId, setLastRecoId] = useState(null);

  // Retrieve last recommendation ID from localStorage
  useEffect(() => {
    const rid = localStorage.getItem("last_reco_id");
    if (rid) setLastRecoId(parseInt(rid, 10));
  }, []);

  // Fetch list of saved summaries
  const fetchSummaries = async () => {
    setLoadingSummaries(true);
    try {
      const res = await api.get("/summaries");
      setSummaries(res.data);
    } catch (err) {
      console.error(err);
    } finally {
      setLoadingSummaries(false);
    }
  };

  useEffect(() => {
    fetchSummaries();
  }, []);

  // 🔹 Fetch AI-generated course recommendations
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

  // Automatically trigger recommendations when tab loads
  useEffect(() => {
    fetchRecommendations();
  }, []);

  // 🔹 Fetch scholarships via web search + AI
  const fetchScholarships = async () => {
    setFetchingScholarships(true);
    try {
      const res = await api.post("/scholarships/generate");
      const list = Array.isArray(res.data?.scholarships)
        ? res.data.scholarships
        : [];
      setScholarships(list);
      if (list.length === 0) {
        alert("No scholarships found for this profile yet.");
      }
    } catch (e) {
      console.error(e);
      alert(e.response?.data?.error || "Failed to fetch scholarships.");
    } finally {
      setFetchingScholarships(false);
    }
  };

  // 🔹 Save summary (PDF)
  const saveSummary = async () => {
    if (!lastRecoId) return alert("No recommendation available to save yet.");
    setSaving(true);
    try {
      await api.post("/summaries", { recommendation_id: lastRecoId });
      alert("✅ Summary PDF saved successfully!");
      await fetchSummaries();
    } catch (e) {
      console.error(e);
      alert(e.response?.data?.error || "Failed to save summary.");
    } finally {
      setSaving(false);
    }
  };

  // 🔹 Download summary
  const handleDownload = async (id) => {
    try {
      const response = await api.get(`/summaries/${id}/download`, {
        responseType: "blob",
      });
      const blob = new Blob([response.data], { type: "application/pdf" });
      const url = window.URL.createObjectURL(blob);
      const link = document.createElement("a");
      link.href = url;
      link.setAttribute("download", `summary_${id}.pdf`);
      document.body.appendChild(link);
      link.click();
      link.parentNode.removeChild(link);
      window.URL.revokeObjectURL(url);
    } catch (error) {
      console.error(error);
      alert("Failed to download PDF");
    }
  };

  // 🔹 Delete summary
  const handleDelete = async (id) => {
    if (!window.confirm("Are you sure you want to delete this summary?")) return;
    try {
      await api.delete(`/summaries/${id}`);
      alert("🗑️ Summary deleted successfully");
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
        <div className="rounded-lg border border-gray-300 bg-white p-6">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm text-gray-500">Documents Analyzed</p>
              <p className="mt-1 text-3xl font-bold text-gray-900">
                {uploadedDocuments.length}
              </p>
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
              <p className="mt-1 text-3xl font-bold text-gray-900">
                {courses.length}
              </p>
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
              <p className="mt-1 text-3xl font-bold text-gray-900">
                {scholarships.length}
              </p>
            </div>
            <div className="rounded-full bg-green-100 p-3">
              <Award className="h-6 w-6 text-green-600" />
            </div>
          </div>
        </div>
      </div>

      {/* Recommended Courses */}
      <div>
        <h2 className="mb-4 text-xl font-bold text-gray-900">
          Recommended Courses
        </h2>
        <div className="grid gap-4">
          {courses.length === 0 && (
            <p className="text-gray-500 text-sm">
              No recommendations yet. Try uploading a transcript.
            </p>
          )}
          {courses.map((course, idx) => (
            <div
              key={idx}
              className="rounded-lg border border-gray-300 bg-white p-6 hover:shadow-md transition-shadow"
            >
              <div className="flex items-start justify-between gap-4">
                <div className="flex-1">
                  <h3 className="font-semibold text-gray-900">
                    {course.title}
                  </h3>
                  <p className="mt-1 text-sm text-gray-500">
                    {course.description}
                  </p>
                </div>
                <div className="flex flex-col items-end gap-2">
                  <div className="rounded-full bg-blue-100 px-3 py-1">
                    <span className="text-sm font-semibold text-blue-600">
                      {Math.round(course.match)}%
                    </span>
                  </div>
                  <button className="text-xs font-medium text-blue-600 hover:underline">
                    Learn More →
                  </button>
                </div>
              </div>
            </div>
          ))}
        </div>
      </div>

      {/* Scholarship finder header + button */}
      <div className="flex items-center justify-between">
        <h2 className="text-xl font-bold text-gray-900">
          Scholarship Opportunities
        </h2>
        <button
          onClick={fetchScholarships}
          disabled={fetchingScholarships}
          className="inline-flex items-center gap-2 rounded-lg bg-green-600 px-4 py-2 text-white font-semibold hover:bg-green-700 disabled:opacity-60"
        >
          <Globe className="w-4 h-4" />
          {fetchingScholarships ? "Searching..." : "Find Scholarships"}
        </button>
      </div>

      {/* Scholarships list */}
      <div className="grid gap-4">
        {fetchingScholarships && (
          <div className="flex items-center justify-center py-8 text-gray-600">
            <Loader2 className="h-5 w-5 animate-spin mr-2 text-green-600" />
            Searching scholarships...
          </div>
        )}

        {!fetchingScholarships && scholarships.length === 0 && (
          <p className="text-gray-500 text-sm">
            No scholarships loaded yet. Click{" "}
            <b>Find Scholarships</b> to search the web based on your transcript.
          </p>
        )}

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
  );
}
