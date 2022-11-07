using DevSmtp.Core.Models;

namespace DevSmtp.Core.Queries
{
    public sealed class FindMessageByIdResult : QueryResult
    {
        public FindMessageByIdResult(Message? message)
        {
            this.Message = message;
        }

        public FindMessageByIdResult(Exception error)
            : base(error)
        {
        }

        public Message? Message { get; }
    }
}
