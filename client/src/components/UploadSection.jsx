import React, { useState } from "react"
import { Upload, FileText, CheckCircle } from "lucide-react"

export default function UploadSection({ onUpload }) {
  const [dragActive, setDragActive] = useState(false)
  const [uploadedFiles, setUploadedFiles] = useState([])

  const handleDrag = (e) => {
    e.preventDefault()
    e.stopPropagation()
    if (e.type === "dragenter" || e.type === "dragover") {
      setDragActive(true)
    } else if (e.type === "dragleave") {
      setDragActive(false)
    }
  }

  const handleDrop = (e) => {
    e.preventDefault()
    e.stopPropagation()
    setDragActive(false)

    const files = Array.from(e.dataTransfer.files)
    const fileNames = files.map((f) => f.name)
    setUploadedFiles((prev) => [...prev, ...fileNames])
  }

  const handleFileInput = (e) => {
    const files = Array.from(e.target.files || [])
    if (onUpload) onUpload(files)
    const fileNames = files.map((f) => f.name)
    setUploadedFiles((prev) => [...prev, ...fileNames])
  }

  const handleAnalyze = () => {
    onUpload(uploadedFiles)
  }

  return (
    <div className="space-y-8">
      {/* Upload Area */}
      <div
        onDragEnter={handleDrag}
        onDragLeave={handleDrag}
        onDragOver={handleDrag}
        onDrop={handleDrop}
        className={`relative rounded-xl border-2 border-dashed p-12 text-center transition-all ${dragActive
          ? "border-blue-600 bg-blue-100"
          : "border-gray-300 bg-gray-100 hover:border-blue-400"
          }`}
      >
        <div className="flex flex-col items-center gap-4">
          <div className={`relative rounded-full bg-blue-100 p-6
            transition-all duration-300 ${dragActive ? "scale-110 bg-blue-200" : "scale-100"}`}>
            <div className="absolute inset-0 rounded-full bg-blue-200 animate-ping" />
            <Upload
              className={`relative h-12 w-12 text-blue-600
                transition-transform duration-300
                ${dragActive ? "-translate-y-1" : "translate-y-0"}`}
            />
          </div>
          <div>
            <h3 className="text-lg font-semibold text-gray-900">
              Upload Your Academic Documents
            </h3>
            <p className="mt-1 text-sm text-gray-500">
              Drag and drop your transcript, course history, or other academic files
            </p>
          </div>
          <label className="cursor-pointer group">
            <input
              type="file"
              multiple
              onChange={handleFileInput}
              className="hidden"
              accept=".pdf,.doc,.docx,.txt"
            />
            <span className="
              inline-flex items-center gap-3 rounded-xl
              bg-blue-600 px-8 py-3.5
              font-semibold text-white
              transition-all duration-200
              hover:bg-blue-500 hover:scale-105 hover:shadow-lg
              active:scale-95
            ">
              <Upload className="h-5 w-5" />
              Browse Files
            </span>
          </label>
        </div>
      </div>

      {/* Uploaded Files */}
      {uploadedFiles.length > 0 && (
        <div className="space-y-4">
          <h3 className="font-semibold text-gray-900">Uploaded Documents</h3>
          <div className="grid gap-3">
            {uploadedFiles.map((file, idx) => (
              <div
                key={idx}
                className="flex items-center gap-3 rounded-lg border border-gray-300 bg-white p-4 hover:bg-gray-100 transition-colors"
              >
                <FileText className="h-5 w-5 text-blue-600" />
                <div className="flex-1">
                  <p className="font-medium text-gray-900">{file}</p>
                  <p className="text-xs text-gray-500">Ready for analysis</p>
                </div>
                <CheckCircle className="h-5 w-5 text-green-500" />
              </div>
            ))}
          </div>

          {/* Analyze Button */}
          <button
            onClick={handleAnalyze}
            className="w-full rounded-lg bg-blue-600 px-6 py-3 font-semibold text-white hover:bg-blue-500 transition-colors mt-6"
          >
            Analyze Documents & Get Recommendations
          </button>
        </div>
      )}

      {/* Info Cards */}
      <div className="grid gap-4 md:grid-cols-3">
        <div className="rounded-lg border border-gray-300 bg-white p-6">
          <div className="mb-3 text-2xl">📚</div>
          <h4 className="font-semibold text-gray-900">Course Matching</h4>
          <p className="mt-2 text-sm text-gray-500">
            Get personalized course recommendations based on your academic profile
          </p>
        </div>
        <div className="rounded-lg border border-gray-300 bg-white p-6">
          <div className="mb-3 text-2xl">🎓</div>
          <h4 className="font-semibold text-gray-900">Scholarship Finder</h4>
          <p className="mt-2 text-sm text-gray-500">
            Discover scholarships that match your qualifications and interests
          </p>
        </div>
        <div className="rounded-lg border border-gray-300 bg-white p-6">
          <div className="mb-3 text-2xl">🤖</div>
          <h4 className="font-semibold text-gray-900">AI Advisor</h4>
          <p className="mt-2 text-sm text-gray-500">
            Chat with our AI to get detailed guidance on your academic journey
          </p>
        </div>
      </div>
    </div>

  )
}



/* new design */

// import React, { useState } from "react"
// import { Upload, FileText, CheckCircle2, X } from "lucide-react"

// export default function UploadSection({ onUpload }) {
//   const [isDragging, setIsDragging] = useState(false)
//   const [files, setFiles] = useState([])

//   const handleFileInput = (e) => {
//     if (e.target.files) {
//       setFiles(Array.from(e.target.files))
//     }
//     if (onUpload) onUpload(Array.from(e.target.files))
//   }

//   const handleDragOver = (e) => {
//     e.preventDefault()
//     setIsDragging(true)
//   }

//   const handleDragLeave = (e) => {
//     e.preventDefault()
//     setIsDragging(false)
//   }

//   const handleDrop = (e) => {
//     e.preventDefault()
//     setIsDragging(false)
//     if (e.dataTransfer.files) {
//       setFiles(Array.from(e.dataTransfer.files))
//     }
//   }

//   const removeFile = (index) => {
//     setFiles(files.filter((_, i) => i !== index))
//   }

//   return (
//     <div >
//       <div className="min-screen flex items-center justify-center bg-gray-100 p-4">
//         <div className="w-full max-w-2xl">
//           <div className="text-center mb-8">
//             <h1 className="text-4xl font-bold text-gray-900 mb-3">
//               Upload Your Academic Documents
//             </h1>
//             <p className="text-lg text-gray-500">
//               Securely upload transcripts, course history, and academic files
//             </p>
//           </div>

//           <div
//             onDragOver={handleDragOver}
//             onDragLeave={handleDragLeave}
//             onDrop={handleDrop}
//             className={`relative overflow-hidden rounded-2xl border-2 border-dashed
//         transition-all duration-300 ease-out
//         ${isDragging
//                 ? "border-blue-600 bg-blue-100 scale-[1.02]"
//                 : "border-gray-300 bg-white hover:border-blue-400 hover:bg-gray-50"
//               }`}
//           >
//             <div className="p-12">
//               <div className="flex flex-col items-center gap-6">
//                 <div className={`relative rounded-full bg-blue-100 p-6
//             transition-all duration-300 ${isDragging ? "scale-110 bg-blue-200" : "scale-100"}`}>
//                   <div className="absolute inset-0 rounded-full bg-blue-200 animate-ping" />
//                   <Upload
//                     className={`relative h-12 w-12 text-blue-600
//                 transition-transform duration-300
//                 ${isDragging ? "-translate-y-1" : "translate-y-0"}`}
//                   />
//                 </div>

//                 <div className="text-center space-y-2">
//                   <h3 className="text-xl font-semibold text-gray-900">
//                     {isDragging ? "Drop your files here" : "Drag & drop your files"}
//                   </h3>
//                   <p className="text-sm text-gray-500 max-w-md">
//                     Support for PDF, DOC, DOCX, and TXT files up to 10MB each
//                   </p>
//                 </div>

//                 <div className="flex items-center gap-4 w-full max-w-xs">
//                   <div className="h-px flex-1 bg-gray-300" />
//                   <span className="text-xs font-medium text-gray-500 uppercase tracking-wider">or</span>
//                   <div className="h-px flex-1 bg-gray-300" />
//                 </div>

//                 <label className="cursor-pointer group">
//                   <input
//                     type="file"
//                     multiple
//                     onChange={handleFileInput}
//                     className="hidden"
//                     accept=".pdf,.doc,.docx,.txt"
//                   />
//                   <span className="
//               inline-flex items-center gap-3 rounded-xl
//               bg-blue-600 px-8 py-3.5
//               font-semibold text-white
//               transition-all duration-200
//               hover:bg-blue-500 hover:scale-105 hover:shadow-lg
//               active:scale-95
//             ">
//                     <Upload className="h-5 w-5" />
//                     Browse Files
//                   </span>
//                 </label>
//               </div>
//             </div>

//             <div className="absolute top-0 right-0 w-32 h-32 bg-blue-100 rounded-full blur-3xl -translate-y-1/2 translate-x-1/2" />
//             <div className="absolute bottom-0 left-0 w-32 h-32 bg-gray-200 rounded-full blur-3xl translate-y-1/2 -translate-x-1/2" />
//           </div>

//           {files.length > 0 && (
//             <div className="mt-6 space-y-3">
//               <div className="flex items-center justify-between">
//                 <h4 className="text-sm font-semibold text-gray-900">
//                   Selected Files ({files.length})
//                 </h4>
//                 <button
//                   onClick={() => setFiles([])}
//                   className="text-xs text-gray-500 hover:text-gray-900 transition-colors"
//                 >
//                   Clear all
//                 </button>
//               </div>
//               <div className="space-y-2">
//                 {files.map((file, index) => (
//                   <div
//                     key={index}
//                     className="
//                 flex items-center gap-3 p-4 rounded-xl
//                 bg-white border border-gray-300
//                 hover:border-blue-400 transition-all duration-200
//                 animate-in slide-in-from-bottom-2
//               "
//                     style={{ animationDelay: `${index * 50}ms` }}
//                   >
//                     <div className="flex-shrink-0 rounded-lg bg-blue-100 p-2">
//                       <FileText className="h-5 w-5 text-blue-600" />
//                     </div>
//                     <div className="flex-1 min-w-0">
//                       <p className="text-sm font-medium text-gray-900 truncate">{file.name}</p>
//                       <p className="text-xs text-gray-500">{(file.size / 1024).toFixed(1)} KB</p>
//                     </div>
//                     <CheckCircle2 className="h-5 w-5 text-green-500 flex-shrink-0" />
//                     <button
//                       onClick={() => removeFile(index)}
//                       className="
//                   flex-shrink-0 rounded-lg p-1.5
//                   text-gray-500 hover:text-gray-900
//                   hover:bg-gray-100 transition-colors
//                 "
//                     >
//                       <X className="h-4 w-4" />
//                     </button>
//                   </div>
//                 ))}
//               </div>
//             </div>
//           )}

//           {files.length > 0 && (
//             <button
//               className="
//           mt-6 w-full rounded-xl bg-blue-600 px-6 py-4
//           font-semibold text-white
//           transition-all duration-200
//           hover:bg-blue-500 hover:shadow-lg
//           active:scale-[0.98]
//         "
//             >
//               Upload {files.length} {files.length === 1 ? "File" : "Files"}
//             </button>
//           )}
//         </div>

//       </div>
//     </div>


//   )
// }


