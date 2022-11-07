using DevSmtp.Core.Models;

namespace DevSmtp.Core.Queries
{
    public sealed class FindMessagesByEmailResult : QueryResult
    {
        public FindMessagesByEmailResult(IEnumerable<Message>? messages)
        {
            this.Messages = messages;
        }

        public FindMessagesByEmailResult(Exception error)
            : base(error)
        {
            this.Messages = new List<Message>();
        }

        public IEnumerable<Message> Messages { get; }
    }
}
