using DevSmtp.Core.Models;

namespace DevSmtp.Core.Queries
{
    public sealed class FindMessageById : IQuery<FindMessageByIdResult>
    {
        public FindMessageById(MessageId id)
        {
            this.Id = id;
        }

        public MessageId Id { get; }
    }
}
