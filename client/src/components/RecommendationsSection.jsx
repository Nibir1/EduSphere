"use client"

import React from "react"
import { TrendingUp, Award, BookOpen } from "lucide-react"

export default function RecommendationsSection({ recommendations, uploadedDocuments }) {
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
          {courses.map((course, idx) => (
            <div key={idx} className="rounded-lg border border-gray-300 bg-white p-6 hover:shadow-md transition-shadow">
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
          {scholarships.map((scholarship, idx) => (
            <div key={idx} className="rounded-lg border border-gray-300 bg-white p-6 hover:shadow-md transition-shadow">
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
    </div>
  )
}
